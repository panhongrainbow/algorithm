package toolset

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// createTestPool creates a new NumberPool with the provided values.
// It initializes a NumberPool and populates it with the given values.
func createTestPool[T Number](values []T) *NumberPool[T] {
	// Create a new NumberPool instance.
	np := NewNumberPool[T]()

	// Iterate over each value in the provided slice.
	for _, v := range values {
		// Add the value to the pool with an empty struct as its value.
		np.pool[v] = struct{}{}
	}

	// Set the keyCount of the pool to the number of values provided.
	np.keyCount = len(values)

	// Return the populated NumberPool instance.
	return np
}

// Test_GenerateUniqueInt64Numbers_WithPool tests the GenerateUniqueNumbers function of NumberPool.
// It uses different test cases to verify the correctness of the function.
func Test_GenerateUniqueInt64Numbers_WithPool(t *testing.T) {
	// Define a slice of test cases with various parameters and expected results.
	tests := []struct {
		name          string  // Name of the test case.
		minNum        int64   // Minimum number in the range.
		maxNum        int64   // Maximum number in the range.
		count         int     // Number of unique numbers to generate.
		withdraw      int     // Number of unique numbers to withdraw from the pool.
		fullRemove    bool    // Whether to remove all numbers from the pool.
		shuffle       bool    // Whether to shuffle the numbers in the pool.
		existingPool  []int64 // Existing pool values to start with.
		expectedError string  // Expected error message, if any.
	}{
		{
			name:          "Withdraw more than possible",             // Test case name.
			minNum:        1,                                         // Minimum number for generation.
			maxNum:        10,                                        // Maximum number for generation.
			count:         2,                                         // Number of unique numbers to generate.
			withdraw:      10,                                        // Number of unique numbers to withdraw.
			fullRemove:    false,                                     // Do not remove all numbers from the pool.
			shuffle:       false,                                     // Do not shuffle the numbers in the pool.
			expectedError: "withdraw amount exceeds final pool size", // Expected error message.
		},
		{
			name:         "Full remove, no shuffle", // Test case name.
			minNum:       1,                         // Minimum number for generation.
			maxNum:       10,                        // Maximum number for generation.
			count:        2,                         // Number of unique numbers to generate.
			withdraw:     5,                         // Number of unique numbers to withdraw.
			fullRemove:   true,                      // Remove all numbers from the pool.
			shuffle:      false,                     // Do not shuffle the numbers in the pool.
			existingPool: []int64{1, 2, 3},          // Existing pool values.
		},
		{
			name:         "Partial remove, no shuffle", // Test case name.
			minNum:       1,                            // Minimum number for generation.
			maxNum:       10,                           // Maximum number for generation.
			count:        2,                            // Number of unique numbers to generate.
			withdraw:     2,                            // Number of unique numbers to withdraw.
			fullRemove:   false,                        // Do not remove all numbers from the pool.
			shuffle:      false,                        // Do not shuffle the numbers in the pool.
			existingPool: []int64{1, 2, 3},             // Existing pool values.
		},
		{
			name:         "Partial remove, shuffle", // Test case name.
			minNum:       1,                         // Minimum number for generation.
			maxNum:       10,                        // Maximum number for generation.
			count:        2,                         // Number of unique numbers to generate.
			withdraw:     2,                         // Number of unique numbers to withdraw.
			fullRemove:   false,                     // Do not remove all numbers from the pool.
			shuffle:      true,                      // Shuffle the numbers in the pool.
			existingPool: []int64{1, 2, 3},          // Existing pool values.
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run each test case with a subtest name.
		t.Run(tt.name, func(t *testing.T) {
			// Create a NumberPool with the existing pool values.
			np := createTestPool[int64](tt.existingPool)

			// Call GenerateUniqueNumbers with the test case parameters.
			newNums, removedNums, err := np.GenerateUniqueNumbers(
				tt.minNum, tt.maxNum,
				WithBasicOpt(tt.count, tt.withdraw, tt.fullRemove),
				WithAdvanceOpt(tt.shuffle),
			)

			// Check if an error is expected.
			if tt.expectedError != "" {
				// Assert that an error occurred.
				require.Error(t, err)
				// Assert that the error message matches the expected error.
				assert.Equal(t, tt.expectedError, err.Error())
				// Skip further checks if an error is expected.
				return
			}

			// Check that no error occurred.
			require.NoError(t, err)

			// Extract sorted keys from the pool.
			poolKeys := np.ExtractSortedKeys()
			// Create a map to track seen keys.
			seen := make(map[int64]struct{})
			// Iterate over the keys in the pool.
			for _, num := range poolKeys {
				// Assert that the key is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in the pool")
				// Mark the key as seen.
				seen[num] = struct{}{}
			}

			// Reset the seen map for new numbers.
			seen = make(map[int64]struct{})
			// Iterate over the newly generated numbers.
			for _, num := range newNums {
				// Assert that the new number is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in newNums")
				// Mark the new number as seen.
				seen[num] = struct{}{}
			}

			// Reset the seen map for removed numbers.
			seen = make(map[int64]struct{})
			// Iterate over the removed numbers.
			for _, num := range removedNums {
				// Assert that the removed number is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in removedNums")
				// Mark the removed number as seen.
				seen[num] = struct{}{}
			}

			// Check that all keys in the pool are within the specified range.
			for _, num := range np.ExtractSortedKeys() {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}

			// Check that all newly generated numbers are within the specified range.
			for _, num := range newNums {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}

			// Check that all removed numbers are within the specified range.
			for _, num := range removedNums {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}
		})
	}
}

// Test_GenerateUniqueFloat64Numbers_WithPool tests GenerateUniqueNumbers function for float64 type.
func Test_GenerateUniqueFloat64Numbers_WithPool(t *testing.T) {
	// Define test cases for float64.
	tests := []struct {
		name          string    // Name of the test case.
		minNum        float64   // Minimum number in the range.
		maxNum        float64   // Maximum number in the range.
		count         int       // Number of unique numbers to generate.
		withdraw      int       // Number of unique numbers to withdraw from the pool.
		fullRemove    bool      // Whether to remove all numbers from the pool.
		shuffle       bool      // Whether to shuffle the numbers in the pool.
		existingPool  []float64 // Existing pool values to start with.
		expectedError string    // Expected error message, if any.
	}{
		{
			name:          "Withdraw more than possible",             // Test case name.
			minNum:        1.5,                                       // Minimum number for generation.
			maxNum:        10.5,                                      // Maximum number for generation.
			count:         2,                                         // Number of unique numbers to generate.
			withdraw:      10,                                        // Number of unique numbers to withdraw.
			fullRemove:    false,                                     // Do not remove all numbers from the pool.
			shuffle:       false,                                     // Do not shuffle the numbers in the pool.
			expectedError: "withdraw amount exceeds final pool size", // Expected error message.
		},
		{
			name:         "Full remove, no shuffle", // Test case name.
			minNum:       1.5,                       // Minimum number for generation.
			maxNum:       10.5,                      // Maximum number for generation.
			count:        2,                         // Number of unique numbers to generate.
			withdraw:     5,                         // Number of unique numbers to withdraw.
			fullRemove:   true,                      // Remove all numbers from the pool.
			shuffle:      false,                     // Do not shuffle the numbers in the pool.
			existingPool: []float64{1.5, 2.5, 3.5},  // Existing pool values.
		},
		{
			name:         "Partial remove, no shuffle", // Test case name.
			minNum:       1.5,                          // Minimum number for generation.
			maxNum:       10.5,                         // Maximum number for generation.
			count:        2,                            // Number of unique numbers to generate.
			withdraw:     2,                            // Number of unique numbers to withdraw.
			fullRemove:   false,                        // Do not remove all numbers from the pool.
			shuffle:      false,                        // Do not shuffle the numbers in the pool.
			existingPool: []float64{1.5, 2.5, 3.5},     // Existing pool values.
		},
		{
			name:         "Partial remove, shuffle", // Test case name.
			minNum:       1.5,                       // Minimum number for generation.
			maxNum:       10.5,                      // Maximum number for generation.
			count:        2,                         // Number of unique numbers to generate.
			withdraw:     2,                         // Number of unique numbers to withdraw.
			fullRemove:   false,                     // Do not remove all numbers from the pool.
			shuffle:      true,                      // Shuffle the numbers in the pool.
			existingPool: []float64{1.5, 2.5, 3.5},  // Existing pool values.
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run each test case with a subtest name.
		t.Run(tt.name, func(t *testing.T) {
			// Create a NumberPool with the existing pool values.
			np := createTestPool[float64](tt.existingPool)

			// Call GenerateUniqueNumbers with the test case parameters.
			newNums, removedNums, err := np.GenerateUniqueNumbers(
				tt.minNum, tt.maxNum,
				WithBasicOpt(tt.count, tt.withdraw, tt.fullRemove),
				WithAdvanceOpt(tt.shuffle),
			)

			// Check if an error is expected.
			if tt.expectedError != "" {
				// Assert that an error occurred.
				require.Error(t, err)
				// Assert that the error message matches the expected error.
				assert.Equal(t, tt.expectedError, err.Error())
				// Skip further checks if an error is expected.
				return
			}

			// Check that no error occurred.
			require.NoError(t, err)

			// Extract sorted keys from the pool.
			poolKeys := np.ExtractSortedKeys()
			// Create a map to track seen keys.
			seen := make(map[float64]struct{})
			// Iterate over the keys in the pool.
			for _, num := range poolKeys {
				// Assert that the key is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in the pool")
				// Mark the key as seen.
				seen[num] = struct{}{}
			}

			// Reset the seen map for new numbers.
			seen = make(map[float64]struct{})
			// Iterate over the newly generated numbers.
			for _, num := range newNums {
				// Assert that the new number is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in newNums")
				// Mark the new number as seen.
				seen[num] = struct{}{}
			}

			// Reset the seen map for removed numbers.
			seen = make(map[float64]struct{})
			// Iterate over the removed numbers.
			for _, num := range removedNums {
				// Assert that the removed number is not a duplicate.
				assert.NotContains(t, seen, num, "Duplicate key found in removedNums")
				// Mark the removed number as seen.
				seen[num] = struct{}{}
			}

			// Check that all keys in the pool are within the specified range.
			for _, num := range np.ExtractSortedKeys() {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}

			// Check that all newly generated numbers are within the specified range.
			for _, num := range newNums {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}

			// Check that all removed numbers are within the specified range.
			for _, num := range removedNums {
				assert.GreaterOrEqual(t, num, tt.minNum)
				assert.LessOrEqual(t, num, tt.maxNum)
			}
		})
	}
}
