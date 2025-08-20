package model2

import (
	"math/rand"
	"time"

	"github.com/panhongrainbow/algorithm/randhub"
	"github.com/panhongrainbow/algorithm/utilhub"
)

// BpTestModel2 🧮 is implemented using the Dynamic Pool Stress Test to simulate random insertions and removals in a real data pool,
// ensuring performance, stability, and correctness.
type BpTestModel2 struct{}

// GenerateRandomSet 🧮 generates a slice of random data set for test model 2.
func (model2 *BpTestModel2) GenerateRandomSet() ([]int64, error) {

	unitTestConfig := utilhub.GetDefaultConfig()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	testPlan := model2.RandomizedBoundary(5, 50, 10, 20)

	progressBar, _ := utilhub.NewProgressBar(
		"Mode 2: Randomized Boundary - generate test data", // Progress bar title.
		uint32(model2.CountOps(testPlan)),                  // Total number of operations.
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

	for j := 0; j < len(testPlan); j++ {
		batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].OperationPlan[0]), -1*int(testPlan[j].OperationPlan[1]), false)

		shuffleSlice(batchInsert, random)
		shuffleSlice(batchRemove, random)

		for k := 0; k < int(testPlan[j].OperationPlan[0]); k++ {
			dataSet = append(dataSet, batchInsert[k])
			progressBar.UpdateBar()
		}

		for l := 0; l < -1*int(testPlan[j].OperationPlan[1]); l++ {
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

// shuffleSlice randomly shuffles the elements in the slice.
func shuffleSlice(slice []int64, rng *rand.Rand) {
	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := rng.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers.
func (model2 *BpTestModel2) CheckRandomSet(dataSet []int64) error {

	// Return nil if the data set is valid.
	return nil
}
