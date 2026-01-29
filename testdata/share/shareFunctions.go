package share

import (
	"errors"
	"math/rand"
	"time"

	"github.com/panhongrainbow/go-algorithm/randhub"
	"github.com/panhongrainbow/go-algorithm/utilhub"
)

type BpTestShare struct{}

// ShareGenerateRandomSet 🧮 generates a slice of random data set for test model 2 and test model 3.
func (model *BpTestShare) ShareGenerateRandomSet(cyclicStressCount int64) ([]int64, error) {
	// Use RandomTotalCount to limit the test scope.
	unitTestConfig := utilhub.GetDefaultConfig()
	limitTestScope := unitTestConfig.Parameters.RandomTotalCount
	stageParams := unitTestConfig.PoolStage

	testPlan := model.StageParameters(limitTestScope, stageParams.MinRemovals, stageParams.MaxRemovals, stageParams.MinPreserveInPool, stageParams.MaxPreserveInPool)

	progressBar, _ := utilhub.NewProgressBar(
		"Mode 3: CyclicStress Boundary - generate test data", // Progress bar title.
		uint32(model._TotalOps(testPlan, cyclicStressCount)), // Total number of operations.
		70,                                      // Progress bar width.
		utilhub.WithTracking(5),                 // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),     // Time zone.
		utilhub.WithTimeControl(500),            // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightBlue), // Display style.
	)

	go func() {
		progressBar.ListenPrinter()
	}()

	pool := randhub.NewDoublePool()

	dataSet := make([]int64, 0)

	for j := 0; j < len(testPlan); j++ {
		batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].op.insertAction), int(testPlan[j].op.deleteAction), false)

		for cycle := 0; cycle < int(cyclicStressCount); cycle++ {

			source := rand.NewSource(time.Now().UnixNano())
			random := rand.New(source)

			ShuffleSlice(batchInsert, random)
			// shuffleSlice(batchRemove, random)

			for k := 0; k < int(testPlan[j].op.insertAction); k++ {
				dataSet = append(dataSet, batchInsert[k])
				progressBar.UpdateBar()
			}

			for l := 0; l < int(testPlan[j].op.insertAction); l++ {
				dataSet = append(dataSet, -1*batchInsert[l])
				progressBar.UpdateBar()
			}
		}

		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		ShuffleSlice(batchInsert, random)
		ShuffleSlice(batchRemove, random)

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

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers for test model 2 and test model 3.
func (model *BpTestShare) CheckRandomSet(dataSet []int64) error {
	// Check if the length of the data set is even.
	if len(dataSet)%2 != 0 {
		return errors.New("dataSet length must be even")
	}

	// ▓▒░ Create a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		"Mode 3: CyclicStress Boundary Test - check test data", // Progress bar title.
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

// TotalOps 🧮 returns the total number of insert/delete operations across all stages.
//
// The sum of all OperationPlans across all stages is defined to be zero (total inserts equal total deletes).
// Therefore, the total number of operations can be calculated as:
//
//	\Sigma Op.InsertAction * 2 * Repeat * 1 (for test mode 2)
//	\Sigma Op.InsertAction * 2 * Repeat * cyclicStressCount (for test mode 3)
//
// where Op.InsertAction is used as the insertion count.
func (model *BpTestShare) _TotalOps(stages []stage, cyclicStressCount int64) int64 {
	var totalOps int64
	for _, each := range stages {
		if each.Repeat > 1 {
			totalOps += each.op.insertAction * int64(each.Repeat) * 2 * cyclicStressCount
		}
	}
	return totalOps
}

// ShuffleSlice randomly shuffles the elements in the slice.
func ShuffleSlice(slice []int64, rng *rand.Rand) {
	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := rng.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}
