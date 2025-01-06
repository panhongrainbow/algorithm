package bptestModel1

import (
	"errors"
	"fmt"
	"github.com/panhongrainbow/algorithm/randhub"
	bptestUtilhub "github.com/panhongrainbow/algorithm/testplan/bptestplan/utilhub"
	"github.com/panhongrainbow/algorithm/utilhub"
)

// BpTestModel1 🧮 represents a test model for B Plus tree testing.
// It emulates a scenario where random numbers are generated and inserted into a B Plus tree and then deleted.
type BpTestModel1 struct {
	RandomTotalCount uint64 // RandomTotalCount is the total number of random numbers to be kept for testing.
}

// GenerateRandomSet 🧮 generates a slice of random data set for test model 1.
func (model1 *BpTestModel1) GenerateRandomSet(
	randomMin uint64, // randomMin is the minimum value for generating random numbers.
	randomHitCollisionPercentage uint64, // randomHitCollisionPercentage is the percentage of random number hit collision in map insert.
) ([]int64, error) {
	// Validate RandomTotalCount to ensure it is not zero.
	if model1.RandomTotalCount == 0 {
		// Return an error if RandomTotalCount is zero.
		return nil, errors.New("BpTestModel1.RandomTotalCount cannot be zero")
	}

	// Calculate the maximum random value based on RandomTotalCount, randomHitCollisionPercentage, and randomMin.
	randomMax, err := bptestUtilhub.CalculateRandomMax(model1.RandomTotalCount, randomHitCollisionPercentage, randomMin)
	if err != nil {
		// Return a wrapped error if CalculateRandomMax fails.
		return nil, fmt.Errorf("failed to calculate random max: %w", err)
	}

	// Ensure randomEvenCount is at least 2 to maintain data integrity.
	randomEvenCount := utilhub.Adjust2Even(int64(model1.RandomTotalCount))
	if randomEvenCount < 2 {
		// Return an error if randomEvenCount is less than 2.
		return nil, fmt.Errorf("randomEvenCount must be at least 2, got: %d", randomEvenCount)
	}

	// Generate a set of unique random numbers using randhub.GenerateUniqueNumbers.
	bulkAdd, err := randhub.GenerateUniqueNumbers(uint64(randomEvenCount/2), int64(randomMin), int64(randomMax))
	if err != nil {
		// Return a wrapped error if GenerateUniqueNumbers fails.
		return nil, fmt.Errorf("failed to generate unique numbers: %w", err)
	}

	// Create a new slice to store the dataset.
	dataSet := make([]int64, randomEvenCount, randomEvenCount)

	// Copy the generated random numbers to the dataset slice.
	copy(dataSet, bulkAdd)

	// Randomize the order of the bulkAdd slice using utilhub.ShuffleSlice.
	utilhub.ShuffleSlice(bulkAdd)

	// Calculate the length of the bulkAdd slice.
	bulkAddLen := len(bulkAdd)

	// Generate negative numbers by multiplying the bulkAdd slice with -1 and append to the dataset.
	for i := 0; i < bulkAddLen; i++ {
		dataSet[bulkAddLen+i] = -1 * bulkAdd[i]
	}

	// Return the generated dataset.
	return dataSet, nil
}
