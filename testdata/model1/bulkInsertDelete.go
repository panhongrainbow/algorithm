package bptestModel1

import (
	"errors"
	"fmt"
	"math"

	"github.com/panhongrainbow/algorithm/costars/slice2tree"
	"github.com/panhongrainbow/algorithm/randhub"
	bptestUtilhub "github.com/panhongrainbow/algorithm/testdata/utilhub"
	"github.com/panhongrainbow/algorithm/utilhub"
)

// BpTestModel1 🧮 represents a test model for B Plus Tree testing.
// It emulates a scenario where random numbers are generated and inserted into a B Plus Tree and then deleted.
type BpTestModel1 struct{}

// GenerateRandomSet 🧮 generates a slice of random data set for test model 1.
func (model1 *BpTestModel1) GenerateRandomSet(
	randomMin uint64, // randomMin is the minimum value for generating random numbers.
	randomHitCollisionPercentage uint64, // randomHitCollisionPercentage is the percentage of random number hit collision in map insert.
) ([]int64, error) {
	// Use RandomTotalCount to limit the test scope.
	unitTestConfig := utilhub.GetDefaultConfig()
	limitTestScope := uint64(unitTestConfig.Parameters.RandomTotalCount)

	// Validate RandomTotalCount to ensure it is not zero.
	// I make sure that RandomTotalCount is not zero to enough data for testing.
	if limitTestScope == 0 {
		// Return an error if RandomTotalCount is zero.
		return nil, errors.New("BpTestModel1.RandomTotalCount cannot be zero")
	}

	// Calculate the maximum random value based on RandomTotalCount, randomHitCollisionPercentage, and randomMin.
	randomMax, err := bptestUtilhub.CalculateRandomMax(limitTestScope, randomHitCollisionPercentage, randomMin)
	if err != nil {
		// Return a wrapped error if CalculateRandomMax fails.
		return nil, fmt.Errorf("failed to calculate random max: %w", err)
	}

	// Ensure randomMax is within valid limits to prevent overflow.
	// This check is necessary to ensure the generated numbers can be represented as int64. (检查是否溢出)
	if randomMax > math.MaxInt64 {
		return nil, fmt.Errorf("randomMax must be less than %d, got: %d", math.MaxInt64, randomMax)
	}

	// Ensure randomEvenCount is at least 2 to maintain data integrity.
	// I make sure that enough data is generated for testing again.
	randomEvenCount := utilhub.Adjust2Even(int64(limitTestScope))
	if randomEvenCount < 2 {
		// Return an error if randomEvenCount is less than 2.
		return nil, fmt.Errorf("randomEvenCount must be at least 2, got: %d", randomEvenCount)
	}

	// Generate a set of unique random numbers using randhub.GenerateUniqueNumbers.
	// Then separating the generated numbers into positive and negative numbers.
	bulkAdd, err := randhub.GenerateUniqueNumbers(uint64(randomEvenCount/2), int64(randomMin), int64(randomMax))
	if err != nil {
		// Return a wrapped error if GenerateUniqueNumbers fails.
		return nil, fmt.Errorf("failed to generate unique numbers: %w", err)
	}

	// Creating a new slice to store the dataset, which will be tested.
	dataSet := make([]int64, randomEvenCount, randomEvenCount)

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		"Mode 1: Bulk Insert/Delete - generate test data", // Progress bar title.
		uint32(randomEvenCount),                           // Total number of operations.
		70,                                                // Progress bar width.
		utilhub.WithTracking(5),                           // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),               // Time zone.
		utilhub.WithTimeControl(500),                      // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightBlue),           // Display style.
	)

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// Copying the generated random numbers, positive ones, to the dataset slice.
	copy(dataSet, bulkAdd)

	// Randomizing the order of the bulkAdd slice using utilhub.ShuffleSlice.
	utilhub.ShuffleSlice(bulkAdd)

	// ▓▒░ Updating the progress bar.
	progressBar.AddSpecificTimes(uint32(randomEvenCount / 2))

	// Calculating the length of the bulkAdd slice.
	bulkAddLen := len(bulkAdd)

	// Generating negative numbers by multiplying the bulkAdd slice with -1 and append to the dataset.
	for i := 0; i < bulkAddLen; i++ {
		// Appending the negative number to the dataset.
		dataSet[bulkAddLen+i] = -1 * bulkAdd[i]
		// ▓▒░ Updating the progress bar.
		progressBar.UpdateBar()
	}

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Return the generated dataset.
	return dataSet, nil
}

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers.
func (model1 *BpTestModel1) CheckRandomSet(dataSet []int64) error {
	// Check if the length of the data set is even.
	if len(dataSet)%2 != 0 {
		return errors.New("dataSet length must be even")
	}

	// Create two heaps to store the positive and negative numbers.
	postiveHeap := slice2tree.NewHeap(len(dataSet) / 2)
	negativeHeap := slice2tree.NewHeap(len(dataSet) / 2)

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		"Mode 1: Bulk Insert/Delete - check test data", // Progress bar title.
		uint32(len(dataSet)/2*3),                       // Total number of operations.
		70,                                             // Progress bar width.
		utilhub.WithTracking(5),                        // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),            // Time zone.
		utilhub.WithTimeControl(500),                   // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightGreen),       // Display style.
	)

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// Iterate over the data set and separate the positive and negative numbers into the heaps.
	for i := 0; i < len(dataSet); i++ {
		switch {
		case dataSet[i] > 0:
			// Push the positive number into the positive heap.
			postiveHeap.Push(dataSet[i])

			// ▓▒░ Updating the progress bar.
			progressBar.UpdateBar()
		case dataSet[i] < 0:
			// Push the negative number into the negative heap.
			negativeHeap.Push(-1 * dataSet[i])

			// ▓▒░ Updating the progress bar.
			progressBar.UpdateBar()
		default:
			// Return an error if the data set contains zero.
			return errors.New("dataSet must not contain 0")
		}
	}

	// Compare the positive and negative numbers in the heaps.
	for i := 0; i < len(dataSet)/2; i++ {
		// Check if the popped numbers from the heaps are equal.
		if postiveHeap.Pop() != negativeHeap.Pop() {
			// Return an error if the numbers are not equal.
			return errors.New("dataSet is not valid")
		}

		// ▓▒░ Updating the progress bar.
		progressBar.UpdateBar()
	}

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Return nil if the data set is valid.
	return nil
}
