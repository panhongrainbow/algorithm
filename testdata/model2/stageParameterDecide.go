package model2

import (
	"math/rand"
	"strconv"
)

// stage 🧮 represents a single phase of the model2 test. (被切割)
// Each stage defines how many records to insert and delete, and may involve reusing previously deleted records. (每阶段都会有 新增 和 删除)
// The stage is repeated according to the specified count. (可重复执行)
type stage struct {
	// StageSummary provides a short description of this test stage.
	StageSummary string

	// Op defines the sequence of operations (e.g., insert/delete counts).
	Op struct {
		InsertAction int64
		DeleteAction int64
	}

	// RepeatCount indicates how many times this stage will be executed.
	Repeat int
}

// totalOps 🧮 returns the total number of insert/delete operations across all stages.
//
// The sum of all OperationPlans across all stages is defined to be zero (total inserts equal total deletes).
// Therefore, the total number of operations can be calculated as:
//
//	\Sigma Op.InsertAction * 2 * Repeat
//
// where Op.InsertAction is used as the insertion count.
func (model2 *BpTestModel2) totalOps(stages []stage) int64 {
	var totalOps int64
	for _, each := range stages {
		if each.Repeat > 1 {
			totalOps += each.Op.InsertAction * int64(each.Repeat) * 2
		}
	}
	return totalOps
}

// StageParameters 🧮 defines the configuration for each stage of the test.

// Parameters:
//   - minRemovals:      minimum number of records to delete per stage
//   - maxRemovals:      maximum number of records to delete per stage
//   - minPreserveInPool: minimum number of records to preserve in the pool
//   - maxPreserveInPool: maximum number of records to preserve in the pool
//
// Returns:
//   - testStages: a list of stages, each containing insertion/deletion counts

func (model2 *BpTestModel2) StageParameters(
	randomTotalCount, minRemovals, maxRemovals, minPreserveInPool, maxPreserveInPool int64) (testStages []stage) {
	// Use RandomTotalCount to limit the test scope.
	limitTestScope := uint64(randomTotalCount)

	// It ensures that the maximum values are strictly greater than the minimum value.
	if minRemovals >= maxRemovals || minPreserveInPool >= maxPreserveInPool {
		panic("max must be greater than min for both removal and insertion ranges")
	}

	// This for-loop continues generating test stages until the accumulated pool size reaches the target total count.
	stageID := 0
	var keepInPool int64 = 0
	for keepInPool < int64(limitTestScope) {
		// removals randomly selects the number of deletions within the range [minRemovals, maxRemovals).
		removals := minRemovals + rand.Int63n(maxRemovals-minRemovals)
		// difference randomly selects the number of records to preserve in the pool within the range [minPreserveInPool, maxPreserveInPool).
		difference := minPreserveInPool + rand.Int63n(maxPreserveInPool-minPreserveInPool)

		// This block constructs a stage that defines how many items will be inserted and deleted.
		testStages = append(testStages, stage{
			StageSummary: "Stage " + strconv.Itoa(stageID),
			Op: struct {
				InsertAction int64
				DeleteAction int64
			}{
				InsertAction: removals + difference,
				DeleteAction: removals,
			},
			Repeat: 1,
		})

		// keepInPool updates the total number of records that remain in the pool after this stage.
		keepInPool += difference
		stageID++
	}

	// Return the list of generated test stages.
	return testStages
}
