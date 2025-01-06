package bptestUtilhub

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestCalculateRandomTotalCount tests the CalculateRandomTotalCount function.
func Test_CalculateRandomTotalCount(t *testing.T) {
	// Define a slice of test cases, each containing a percentage and an expected error.
	tests := []struct {
		percentage  uint64 // The percentage of memory to allocate.
		expectedErr error  // The expected error for the test case.
	}{
		{2, nil}, // Test valid case with 50% memory allocation.
		{0, nil}, // Test valid case with 0% (no allocation).
		{150, fmt.Errorf("invalid percentage: 150. Must be between 0 and 100")}, // Invalid percentage.
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run a sub-test for each test case, with a descriptive name.
		t.Run(fmt.Sprintf("Percentage %d", tt.percentage), func(t *testing.T) {
			// Call the CalculateRandomTotalCount function with the test case's percentage.
			actualCount, actualErr := CalculateRandomTotalCount(tt.percentage)

			// Check if an error is expected for this test case.
			if tt.expectedErr != nil {
				// If an error is expected, verify that an error occurred.
				require.Error(t, actualErr)
				// Verify that the error message matches the expected error message.
				assert.EqualError(t, actualErr, tt.expectedErr.Error())
			} else {
				// If no error is expected, verify that no error occurred.
				require.NoError(t, actualErr)
			}

			// Verify that the calculated count is non-negative.
			assert.GreaterOrEqual(t, actualCount, uint64(0), "calculated count should be non-negative")
		})
	}
}

// Test_CalculateRandomMax tests the CalculateRandomMax function.
func Test_CalculateRandomMax(t *testing.T) {
	// Define a slice of test cases, each containing input parameters and expected results.
	tests := []struct {
		name                         string
		randomTotalCount             uint64 // The total number of random numbers to be generated.
		randomHitCollisionPercentage uint64 // The percentage of random number hit collision in map insert.
		randomMin                    uint64 // The minimum value for generating random numbers.
		expectedRandomMax            uint64 // The expected maximum random value.
		expectedError                error  // The expected error for the test case.
	}{
		{
			name:                         "valid inputs",
			randomTotalCount:             100,
			randomHitCollisionPercentage: 10,
			randomMin:                    1,
			expectedRandomMax:            1001,
			expectedError:                nil,
		},
		{
			name:                         "zero hit collision percentage",
			randomTotalCount:             100,
			randomHitCollisionPercentage: 0,
			randomMin:                    1,
			expectedRandomMax:            0,
			expectedError:                errors.New("randomHitCollisionPercentage cannot be zero"),
		},
		{
			name:                         "large inputs",
			randomTotalCount:             1000000,
			randomHitCollisionPercentage: 50,
			randomMin:                    1000,
			expectedRandomMax:            2001000,
			expectedError:                nil,
		},
		{
			// randomTotalCount must be greater than 100 to avoid randomMax is equal to randomMin.
			// Please refer to CalculateRandomMax inside for more details.
			name:                         "randomTotalCount is less than 100",
			randomTotalCount:             99,
			randomHitCollisionPercentage: 10,
			randomMin:                    1,
			expectedRandomMax:            0,
			expectedError:                errors.New("randomTotalCount cannot be less than 100"),
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run a sub-test for each test case, with a descriptive name.
		t.Run(tt.name, func(t *testing.T) {
			// Call the CalculateRandomMax function with the test case's input parameters.
			actualRandomMax, actualError := CalculateRandomMax(tt.randomTotalCount, tt.randomHitCollisionPercentage, tt.randomMin)

			// Check if an error is expected for this test case.
			if tt.expectedError != nil {
				// If an error is expected, verify that an error occurred.
				require.Error(t, actualError)
				// Verify that the error message matches the expected error message.
				assert.EqualError(t, actualError, tt.expectedError.Error())
			} else {
				// If no error is expected, verify that no error occurred.
				require.NoError(t, actualError)
			}

			// Verify that the calculated maximum random value matches the expected value.
			assert.Equal(t, tt.expectedRandomMax, actualRandomMax)
		})
	}
}
