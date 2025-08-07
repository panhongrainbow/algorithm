package bptestUtilhub

import (
	"errors"
	"fmt"

	"github.com/panhongrainbow/algorithm/utilhub"
)

// CalculateRandomTotalCount 🧮 calculates the total count of random numbers that can be generated
// based on the specified percentage of available memory.
func CalculateRandomTotalCount(
	randomMemoryUsagePercentage uint64, // randomMemoryUsagePercentage is the percentage of memory used for keeping random numbers.
) (uint64, error) {
	// Check if the specified percentage is valid (between 0 and 100).
	if randomMemoryUsagePercentage > 100 {
		// If the percentage is invalid, return an error.
		return 0, fmt.Errorf("invalid percentage: %d. Must be between 0 and 100", randomMemoryUsagePercentage)
	}

	// Attempt to calculate the spare slice size using the utilhub.SpareSliceSize function.
	randomTotalCount, err := utilhub.SpareSliceSize(randomMemoryUsagePercentage)
	if err != nil {
		// If an error occurs during calculation, return a wrapped error with a descriptive message.
		return 0, fmt.Errorf("failed to calculate spare slice size: %w", err)
	}

	// If the calculation is successful, return the random total count and a nil error.
	return randomTotalCount, nil
}

// CalculateRandomMax 🧮 computes randomMax based on the inputs and returns an error if inputs are invalid.
func CalculateRandomMax(
	randomTotalCount uint64, // randomTotalCount is the total number of random numbers to be generated.
	randomHitCollisionPercentage uint64, // randomHitCollisionPercentage is the percentage of random number hit collision in map insert.
	randomMin uint64, // randomMin is the minimum value for generating random numbers.
) (uint64, error) {
	// Check if the hit collision percentage is zero to avoid division by zero error.
	if randomHitCollisionPercentage == 0 {
		// Return an error if the hit collision percentage is zero.
		return 0, errors.New("randomHitCollisionPercentage cannot be zero")
	}

	// Check if the random total count is less than 100 to avoid randomMax is equal to randomMin.
	if randomTotalCount < 100 {
		return 0, errors.New("randomTotalCount cannot be less than 100")
	}

	// Calculate the maximum random value.
	// randomTotalCount must be greater than randomHitCollisionPercentage (100) to avoid randomMax is equal to randomMin.
	// In other words, randomMax must be greater than 100.
	randomMax := randomTotalCount/randomHitCollisionPercentage*100 + randomMin

	// Return the calculated maximum random value and no error.
	return randomMax, nil
}
