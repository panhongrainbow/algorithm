package model2

import (
	"errors"
	"math/rand"
	"time"

	"github.com/panhongrainbow/algorithm/randhub"
	"github.com/panhongrainbow/algorithm/testdata/share"
	"github.com/panhongrainbow/algorithm/utilhub"
)

// BpTestModel2 🧮 is implemented using the Dynamic Pool Stress Test to simulate random insertions and removals in a real data pool,
// ensuring performance, stability, and correctness.
type BpTestModel2 struct{}

// GenerateRandomSet 🧮 generates a slice of random data set for test model 2.
func (model2 *BpTestModel2) GenerateRandomSet() ([]int64, error) {
	// Use RandomTotalCount to limit the test scope.
	unitTestConfig := utilhub.GetDefaultConfig()
	limitTestScope := unitTestConfig.Parameters.RandomTotalCount
	stageParams := unitTestConfig.PoolStage

	testPlan := model2.StageParameters(limitTestScope, stageParams.MinRemovals, stageParams.MaxRemovals, stageParams.MinPreserveInPool, stageParams.MaxPreserveInPool)

	progressBar, _ := utilhub.NewProgressBar(
		"Mode 2: Randomized Boundary - generate test data", // Progress bar title.
		uint32(model2.TotalOps(testPlan)),                  // Total number of operations.
		70,                                                 // Progress bar width.
		utilhub.WithTracking(5),                            // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),                // Time zone.
		utilhub.WithTimeControl(500),                       // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightBlue),            // Display style.
	)

	go func() {
		progressBar.ListenPrinter()
	}()

	pool := randhub.NewDoublePool()

	dataSet := make([]int64, 0)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for j := 0; j < len(testPlan); j++ {
		batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].op.insertAction), int(testPlan[j].op.deleteAction), false)

		share.ShuffleSlice(batchInsert, random)
		share.ShuffleSlice(batchRemove, random)

		for k := 0; k < int(testPlan[j].op.insertAction); k++ {
			dataSet = append(dataSet, batchInsert[k])
			progressBar.UpdateBar()
		}

		for l := 0; l < int(testPlan[j].op.deleteAction); l++ {
			dataSet = append(dataSet, -1*batchRemove[l])
			progressBar.UpdateBar()
		}
	}

	_, removeAll := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, 0, 0, true)
	for m := 0; m < len(removeAll); m++ {
		dataSet = append(dataSet, -1*removeAll[m])
		progressBar.UpdateBar()
	}

	progressBar.Complete()

	<-progressBar.WaitForPrinterStop()

	return dataSet, nil
}

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers.
func (model2 *BpTestModel2) CheckRandomSet(dataSet []int64) error {
	// Check if the length of the data set is even.
	if len(dataSet)%2 != 0 {
		return errors.New("dataSet length must be even")
	}

	// ▓▒░ Create a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		"Mode 2: Randomized Boundary Test - check test data", // Progress bar title.
		uint32(len(dataSet)),                     // Total number of operations.
		70,                                       // Progress bar width.
		utilhub.WithTracking(5),                  // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),      // Time zone.
		utilhub.WithTimeControl(500),             // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightGreen), // Display style.
	)

	// Create an empty map for checking dataSet.
	checkPool := make(map[int64]struct{})

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// Iterate through each element in the dataSet.
	for i := 0; i < len(dataSet); i++ {
		switch {
		case dataSet[i] > 0:
			// Check if the positive number already exists in the checkPool.
			_, exists := checkPool[dataSet[i]]
			if !exists {
				// If it doesn't exist, add it to the checkPool.
				checkPool[dataSet[i]] = struct{}{}
			} else {
				// If it already exists, return an error.
				return errors.New("dataSet is not valid")
			}

			// ▓▒░ Updating the progress bar.
			progressBar.UpdateBar()
		case dataSet[i] < 0:
			// Check if the corresponding positive number exists in the checkPool.
			_, exists := checkPool[-1*dataSet[i]]
			if exists {
				// If it exists, remove it from the checkPool.
				delete(checkPool, -1*dataSet[i])
			} else {
				// If it doesn't exist, return an error.
				return errors.New("dataSet is not valid")
			}

			// ▓▒░ Updating the progress bar.
			progressBar.UpdateBar()
		default:
			// Return an error if the data set contains zero.
			return errors.New("dataSet must not contain 0")
		}
	}

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Return nil if the data set is valid.
	return nil
}
