package testToolset

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

// ===================================================
//                  ⚗️ Stability Testing (NumberPool)
// ===================================================
// 🧪 Stability testing assesses the program’s performance over extended periods.
// 🧪 By continuously adding and deleting data, it evaluates the program’s stability and reliability.
// 🧪 These types of testing helps identify potential issues such as memory leaks, performance degradation, or other problems that may arise during long-term operation.

// NumberPool 🧫 represents a pool of unique numbers with a generic type T.
type NumberPool[T Number] struct {
	pool     map[T]struct{} // Stores the unique numbers in the pool.
	keyCount int            // Tracks the number of keys in the pool.
}

// NewNumberPool 🧫 creates and returns a new NumberPool instance for the type T.
// It initializes the pool as an empty map.
func NewNumberPool[T Number]() *NumberPool[T] {
	return &NumberPool[T]{
		pool:     make(map[T]struct{}),
		keyCount: 0,
	}
}

// ExtractSortedKeys 🧫 returns all keys in the pool as a sorted slice.
// It extracts the keys, sorts them, and then returns the sorted slice.
func (np *NumberPool[T]) ExtractSortedKeys() []T {
	// Check if the pool is empty.
	if np.keyCount == 0 {
		return []T{}
	}

	// Initialize a slice to store the keys.
	keys := make([]T, 0, len(np.pool))

	// Append each key from the pool to the slice.
	for key := range np.pool {
		keys = append(keys, key)
	}

	// Sort the slice of keys.
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// Return the sorted slice of keys.
	return keys
}

// ---------------------------------------------------
//                  ⚗️ Functional Options Pattern
// ---------------------------------------------------
// The Functional Options Pattern is used here to provide a flexible way to configure the behavior of a `npSet` instance.
// This pattern allows for easy extension and modification of configuration options without altering the function signatures.

// npSet represents a set of configuration options for number generation.
type npSet struct {
	count      int  // The number of unique numbers to generate.
	withdraw   int  // The number of unique numbers to withdraw from the pool.
	fullRemove bool // If true, all generated numbers are removed from the pool.
	shuffle    bool // If true, shuffle the order of the generated numbers.
}

// NpOpt is a function type that applies options to a npSet instance.
type NpOpt func(*npSet)

// WithBasicOpt creates an option function that sets basic configuration options.
func WithBasicOpt(count, withdraw int, fullRemove bool) NpOpt {
	return func(s *npSet) {
		s.count = count
		s.withdraw = withdraw
		s.fullRemove = fullRemove
	}
}

// WithAdvanceOpt creates an option function that sets the shuffle configuration option.
func WithAdvanceOpt(shuffle bool) NpOpt {
	return func(s *npSet) {
		s.shuffle = shuffle
	}
}

// newNpSet creates a new instance of npSet with default values and applies provided options.
func newNpSet(opts ...NpOpt) *npSet {
	npset := &npSet{
		count:      0,     // Initialize with default value for the number of unique numbers to generate.
		withdraw:   0,     // Initialize with default value for the number of unique numbers to withdraw from the pool.
		fullRemove: false, // Initialize with default value for whether to remove all generated numbers from the pool.
		shuffle:    false, // Initialize with default value for shuffling the order of generated numbers.
	}
	for _, opt := range opts {
		opt(npset) // Apply each provided option to the npSet instance.
	}
	return npset // Return the configured npSet instance.
}

// GenerateUniqueNumbers 🧫 generates a slice of unique numbers of type T within the specified range and optionally withdraws a number of them from the pool.
//
// Parameters:
// minNum: the minimum value (inclusive) for the range of numbers to generate.
// maxNum: the maximum value (inclusive) for the range of numbers to generate.
// count: the number of unique numbers to generate. The function will attempt to generate this many numbers within the specified range.
// opts: optional functional options to customize the behavior of the function. The options include:
//  - withdraw: the number of unique numbers to withdraw from the pool after generation. If not specified, no numbers are withdrawn.
//  - fullRemove: if true, all generated numbers are removed from the pool after generation, regardless of the `withdraw` option. Default is false.
//  - shuffle: if true, the order of the generated numbers is shuffled before returning. Default is false.
//
// Returns:
// A slice of unique numbers that were generated.
// A slice of numbers that were withdrawn from the pool (if applicable).
// An error, if there is an issue with generating or withdrawing numbers (e.g., if the range is invalid or there are not enough unique numbers available).
//
// Note:
// If `withdraw` option is provided, it removes the specified number of these generated numbers from the pool.
// If `fullRemove` is set to true, it will remove all generated numbers from the pool, ensuring that none of the newly generated numbers remain in the pool.
// If `shuffle` is true, the order of the generated numbers is shuffled before returning them.

func (np *NumberPool[T]) GenerateUniqueNumbers(minNum, maxNum T, opts ...NpOpt) ([]T, []T, error) {
	// Create a new set of options.
	npset := newNpSet(opts...)

	// Check if the requested withdrawal amount exceeds the final pool size.
	if npset.withdraw > len(np.pool)+npset.count {
		return []T{}, []T{}, errors.New("withdraw amount exceeds final pool size")
	}

	// Initialize a slice to store newly generated unique numbers.
	newNumbers := make([]T, 0, npset.count)

	// Initialize a slice to store the numbers that will be removed from the pool.
	removedNumbers := make([]T, 0, npset.withdraw)

	// Create a new random number generator with a seed based on the current time.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate new unique numbers within the range [minNum, maxNum].
	for len(newNumbers) < npset.count {
		num := generateRandomNumber(minNum, maxNum, r)

		// Check if the generated number is already in the pool.
		if _, exists := np.pool[num]; !exists {
			// If not, add it to the pool and the list of new numbers.
			np.pool[num] = struct{}{}
			newNumbers = append(newNumbers, num)
		}
	}

	// If fullRemove is true, remove all numbers from the pool.
	if npset.fullRemove {
		// Convert the keys of the map to a slice.
		keys := make([]T, 0, len(np.pool))
		for num := range np.pool {
			keys = append(keys, num)
		}

		// Optionally shuffle the keys to randomize the order.
		if npset.shuffle {
			rand.Shuffle(len(keys), func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})
		}

		// Remove each number from the pool in the randomized order.
		for _, num := range keys {
			removedNumbers = append(removedNumbers, num)
		}
		// Reset the pool by creating a new empty map.
		np.pool = make(map[T]struct{})
	} else {
		// Convert the keys of the map to a slice.
		keys := make([]T, 0, len(np.pool))
		for num := range np.pool {
			keys = append(keys, num)
		}

		// Optionally shuffle the keys to randomize the order.
		if npset.shuffle {
			rand.Shuffle(len(keys), func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
			})
		}

		// Remove a specific number of elements as specified by the withdraw parameter.
		for _, num := range keys {
			delete(np.pool, num)
			removedNumbers = append(removedNumbers, num)

			// If the required number of withdrawals is met, return the results.
			if len(removedNumbers) >= npset.withdraw {
				return newNumbers, removedNumbers, nil
			}
		}
	}

	// Return the lists of new and removed numbers.
	return newNumbers, removedNumbers, nil
}

// generateRandomNumber 🧫 generates a random number within the specified range for type T.
// This uses type assertions to handle different numeric types.
func generateRandomNumber[T Number](min, max T, r *rand.Rand) T {
	// Convert minNum and maxNum to float64 for comparison
	minFloat := float64(min)
	maxFloat := float64(max)

	return T(minFloat + (maxFloat-minFloat)*r.Float64())

}
