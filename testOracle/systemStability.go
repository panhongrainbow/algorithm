package testOracle

import (
	"math/rand"
	"time"
)

type NumberPool struct {
	pool map[int64]struct{}
}

// NewNumberPool creates a new NumberPool.
func NewNumberPool() *NumberPool {
	return &NumberPool{
		pool: make(map[int64]struct{}),
	}
}

// GenerateNumbers generates `count` new numbers in the pool within the range [min, max]
// and removes `withdraw` numbers randomly from the pool or clears the pool based on `fullRemove`.
func (np *NumberPool) GenerateNumbers(min, max int64, count, withdraw int, fullRemove bool) ([]int64, []int64) {
	newNumbers := make([]int64, 0, count)
	removedNumbers := make([]int64, 0, withdraw)

	// Create a new random number generator with a seed based on the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate new unique numbers within the range [min, max]
	for len(newNumbers) < count {
		num := min + r.Int63n(max-min+1)
		if _, exists := np.pool[num]; !exists {
			np.pool[num] = struct{}{}
			newNumbers = append(newNumbers, num)
		}
	}

	if fullRemove {
		// Remove all numbers from the pool if fullRemove is true
		for num := range np.pool {
			removedNumbers = append(removedNumbers, num)
		}
		np.pool = make(map[int64]struct{}) // Reset the pool to empty
	} else {
		// Remove numbers randomly from the pool if fullRemove is false
		// for len(removedNumbers) < withdraw && len(np.pool) > 0 {
		// Convert the keys of the map to a slice
		// var keys []int64
		for num := range np.pool {
			delete(np.pool, num)
			removedNumbers = append(removedNumbers, num)
			if len(removedNumbers) >= withdraw {
				return newNumbers, removedNumbers
			}
		}

		// Randomly select a key to remove
		/*if len(keys) > 0 {
			index := r.Intn(len(keys))
			num := keys[index]
			delete(np.pool, num)
			removedNumbers = append(removedNumbers, num)
		}*/
		// }
	}

	return newNumbers, removedNumbers
}
