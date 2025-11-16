package model2

import (
	"math/rand"
	"strconv"
)

// stage 🧮 represents a single phase of the model2 test. (被切割成很多阶段)
// Each stage defines how many records to insert and delete, and may involve reusing previously deleted records. (每阶段都会有 新增 和 删除)
// The stage is repeated according to the specified count. (重复执行)
type stage struct {
	// StageSummary provides a short description of this test stage.
	stageSummary string

	// Op defines the sequence of operations (e.g., insert/delete counts).
	op struct {
		insertAction int64
		deleteAction int64
	}

	// RepeatCount indicates how many times this stage will be executed.
	Repeat int
}

// TotalOps 🧮 returns the total number of insert/delete operations across all stages.
//
// The sum of all OperationPlans across all stages is defined to be zero (total inserts equal total deletes).
// Therefore, the total number of operations can be calculated as:
//
//	\Sigma Op.InsertAction * 2 * Repeat
//
// where Op.InsertAction is used as the insertion count.
func (model2 *BpTestModel2) TotalOps(stages []stage) int64 {
	var totalOps int64
	for _, each := range stages {
		if each.Repeat > 1 {
			totalOps += each.op.insertAction * int64(each.Repeat) * 2
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
//
// (这里会决定每个阶段的设定细节)

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
			stageSummary: "Stage " + strconv.Itoa(stageID),
			op: struct {
				insertAction int64
				deleteAction int64
			}{
				insertAction: removals + difference,
				deleteAction: removals,
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
