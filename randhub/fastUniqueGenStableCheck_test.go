package randhub

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

// Test_FastGenUniqueNumbers ensures that `GenerateUniqueNumbers` produces unique values within a given range and
// successfully removes all pool values when `fullRemove` is set to true.
func Test_FastGenUniqueNumbers(t *testing.T) {
	// Test case 1: Check that numbers are within range and unique.
	t.Run("Check unique numbers in range", func(t *testing.T) {
		// Define test parameters for generating unique numbers.
		var minNum int64 = 1   // Minimum value for the number generation range.
		var maxNum int64 = 100 // Maximum value for the number generation range.
		count := 10            // Number of unique numbers to generate.
		withdraw := 5          // Number of numbers to withdraw from the pool.
		fullRemove := false    // Flag indicating whether to fully remove items from the pool.

		// Create a new instance of DoublePool.
		np := NewDoublePool()

		// Call GenerateUniqueNumbers to generate and withdraw numbers from the pool.
		newNumbers, removedNumbers := np.GenerateUniqueInt64Numbers(minNum, maxNum, count, withdraw, fullRemove)

		// Ensure newNumbers slice has the expected length.
		require.Equal(t, count, len(newNumbers), "newNumbers slice should have %d elements", count)

		// Ensure removedNumbers slice has the expected length.
		require.Equal(t, withdraw, len(removedNumbers), "removedNumbers slice should have %d elements", withdraw)

		// Create a map to track unique numbers in newNumbers.
		seenNew := make(map[int64]struct{})

		// Iterate through the newNumbers slice.
		for _, num := range newNumbers {
			// Check if the number is within the specified range.
			assert.GreaterOrEqual(t, num, minNum, "newNumbers value %d is less than min %d", num, minNum)
			assert.LessOrEqual(t, num, maxNum, "newNumbers value %d is greater than max %d", num, maxNum)

			// Ensure the number has not been seen before (unique check).
			_, exists := seenNew[num]
			assert.False(t, exists, "newNumbers contains duplicate value %d", num)

			// Mark the number as seen.
			seenNew[num] = struct{}{}
		}

		// Create a map to track unique numbers in removedNumbers.
		seenRemoved := make(map[int64]struct{})

		// Iterate through the removedNumbers slice.
		for _, num := range removedNumbers {
			// Check if the number is within the specified range.
			assert.GreaterOrEqual(t, num, minNum, "removedNumbers value %d is less than min %d", num, minNum)
			assert.LessOrEqual(t, num, maxNum, "removedNumbers value %d is greater than max %d", num, maxNum)

			// Ensure the number has not been seen before (unique check).
			_, exists := seenRemoved[num]
			assert.False(t, exists, "removedNumbers contains duplicate value %d", num)

			// Mark the number as seen.
			seenRemoved[num] = struct{}{}
		}
	})

	// Test case 2: Check that fullRemove extracts all pool values.
	t.Run("Check full removal from pool", func(t *testing.T) {
		// Define test parameters for generating unique numbers with fullRemove enabled.
		var minNum int64 = 1  // Minimum value for the number generation range.
		var maxNum int64 = 10 // Maximum value for the number generation range.
		count := 0            // Number of unique numbers to generate.
		withdraw := 5         // Number of numbers to withdraw from the pool (irrelevant here as fullRemove is true).
		fullRemove := true    // Set to true to remove all values from the pool.

		// Create a new instance of DoublePool and populate it with existing values.
		np := NewDoublePool()
		np.pool[1] = struct{}{}
		np.pool[2] = struct{}{}
		np.pool[3] = struct{}{}
		np.pool[4] = struct{}{}
		np.pool[5] = struct{}{}

		// Call GenerateUniqueNumbers to generate and fully remove all numbers from the pool.
		_, removedNumbers := np.GenerateUniqueInt64Numbers(minNum, maxNum, count, withdraw, fullRemove)

		// Ensure removedNumbers slice contains all original values from the pool.
		expectedRemoved := []int64{1, 2, 3, 4, 5}
		sort.Slice(removedNumbers, func(i, j int) bool { return removedNumbers[i] < removedNumbers[j] })
		sort.Slice(expectedRemoved, func(i, j int) bool { return expectedRemoved[i] < expectedRemoved[j] })

		assert.Equal(t, expectedRemoved, removedNumbers, "removedNumbers should contain all original pool values")

		// Ensure the pool is empty after full removal.
		assert.Empty(t, np.pool, "The pool should be empty after full removal")
	})
}
