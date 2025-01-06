package utilhub

import (
	"math/rand"
	"time"
)

// =====================================================================================================================
//                  🛠️ Randomizer (Tool)
// Randomizer is a tool for generating random numbers and shuffling slices.
// =====================================================================================================================

// ShuffleSlice ⛏️ randomly shuffles the elements in the slice.
func ShuffleSlice(slice []int64) {

	// Initialize a random number generator.
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := random.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}
