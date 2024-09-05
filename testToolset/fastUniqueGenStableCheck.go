package testToolset

import (
	"math/rand"
	"time"
)

// =====================================================
//                  ⚗️ Quick & Simple Testing (FastPool)
// =====================================================
// 🧪 This function, `GenerateUniqueNumbers`, is designed for fast and straightforward testing.
// 🧪 Its simplicity allows you to generate unique random numbers rapidly
// within a specified range while efficiently removing values from a pool.
// 🧪 By focusing on ease of use and speed, this function is ideal for testing scenarios
// that require minimal setup and instant feedback on correctness.

// FastPool 🧫 defines a structure with a pool of unique int64 numbers.
type FastPool struct {
	pool map[int64]struct{}
}

// NewDoublePool 🧫 initializes and returns a new DoublePool.
func NewDoublePool() *FastPool {
	// Create a new DoublePool with an empty map for storing unique numbers.
	return &FastPool{
		pool: make(map[int64]struct{}),
	}
}

// GenerateUniqueInt64Numbers 🧫 generates a set of unique numbers within a range, adds them to the pool,
// and optionally removes numbers from the pool.
func (np *FastPool) GenerateUniqueInt64Numbers(min, max int64, count, withdraw int, fullRemove bool) ([]int64, []int64) {
	// Create a slice to store the newly generated numbers with an initial capacity of 'count'.
	newNumbers := make([]int64, 0, count)
	// Create a slice to store the removed numbers with an initial capacity of 'withdraw'.
	removedNumbers := make([]int64, 0, withdraw)

	// Initialize a new random number generator with the current time as the seed.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Keep generating numbers until the 'count' of unique numbers is reached.
	for len(newNumbers) < count {
		// Generate a random number within the range [min, max].
		num := min + r.Int63n(max-min+1)
		// Check if the number already exists in the pool.
		if _, exists := np.pool[num]; !exists {
			// If the number is not in the pool, add it.
			np.pool[num] = struct{}{}
			// Append the number to the newNumbers slice.
			newNumbers = append(newNumbers, num)
		}
	}

	// If fullRemove is true, all numbers in the pool will be removed.
	if fullRemove {
		// Iterate through the pool to remove all numbers.
		for num := range np.pool {
			// Add each number to the removedNumbers slice.
			removedNumbers = append(removedNumbers, num)
		}
		// Reset the pool to an empty map after removing all numbers.
		np.pool = make(map[int64]struct{}) // Clear the pool.
	} else {
		// If fullRemove is false, only remove 'withdraw' number of items from the pool.
		for num := range np.pool {
			// Remove the number from the pool.
			delete(np.pool, num)
			// Add the removed number to the removedNumbers slice.
			removedNumbers = append(removedNumbers, num)
			// If we have removed enough numbers (withdraw amount), return the result.
			if len(removedNumbers) >= withdraw {
				// Return the new and removed numbers.
				return newNumbers, removedNumbers
			}
		}
	}

	// Return the newly generated numbers and the removed numbers.
	return newNumbers, removedNumbers
}
