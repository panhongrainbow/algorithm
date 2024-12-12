package bpTree

import (
	"github.com/panhongrainbow/algorithm/randhub"
	"math/rand"
	"time"
)

const (
	random_machine_action_fetch_none = iota + 1000
	random_machine_action_fetch_all
)

type UniqueRandomMachine struct {
	//
}

// UniqueRandomMachine generates a list of unique numbers for bulk insertion and deletion.
func (UniqueRandomMachine) Generate(action int, expectedInsertCount int64, expectedDeleteCount int64) (insertedKeys []int64, deletedKeys []int64, err error) {
	// Generate a list of unique numbers for bulk insertion.
	// This function will panic if an error occurs during number generation.
	bulkAdd, err := randhub.GenerateUniqueNumbers(expectedInsertCount, randomMin, randomMax)
	if err != nil {
		// Panic if an error occurs during number generation.
		panic(err)
	}

	// Initialize a random number generator with the current time as the seed.
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Create a copy of the bulk insertion list and shuffle it for deletion.
	// This is done to simulate random deletion of keys.
	bulkDel := make([]int64, expectedDeleteCount)
	copy(bulkDel, bulkAdd)
	shuffleSlice(bulkDel, random)

	return bulkAdd, bulkDel, nil
}

type UniqueRandomInPoolMachine struct {
	//
}

// UniqueRandomInPoolMachine generates a list of unique numbers for bulk insertion and deletion using a pool.
func (UniqueRandomInPoolMachine) Generate(action int, expectedInsertCount int64, expectedDeleteCount int64) (insertedKeys []int64, deletedKeys []int64, err error) {
	// Create a new double pool for generating unique numbers.
	pool := randhub.NewDoublePool()

	// Determine if we need to fetch all numbers from the pool.
	fetchAll := false
	if action == random_machine_action_fetch_all {
		fetchAll = true
	}

	// Generate random numbers for bulk insertion and deletion.
	// The pool will ensure that the numbers are unique.
	batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(
		randomMin,                        // Minimum random value.
		randomMax,                        // Maximum random value.
		int(expectedInsertCount),         // Expected number of insertions.
		intAbs(int(expectedDeleteCount)), // Expected number of deletions.
		fetchAll,                         // Whether to fetch all numbers from the pool.
	)

	// Initialize a random number generator with the current time as the seed.
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Shuffle the bulk insertion and deletion lists to simulate random order.
	shuffleSlice(batchInsert, random)
	shuffleSlice(batchRemove, random)

	// Return the generated numbers and no error.
	return batchInsert, batchRemove, nil
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
