package model2

import (
	"math/rand"
	"strconv"

	"github.com/panhongrainbow/algorithm/utilhub"
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

// RandomizedBoundary 🧮 plans repeated insertion and deletion of data in the boundary of B Plus tree.
func (model2 *BpTestModel2) RandomizedBoundary(minRemovals, maxRemovals, minDifference, maxDifference int64) (testStages []stage) {
	// 🧪 Create a config instance for B plus tree unit testing and parse default values.
	unitTestConfig := utilhub.GetDefaultConfig()
	randomTotalCount := uint64(unitTestConfig.Parameters.RandomTotalCount)

	// Perform boundary check outside the loop to ensure valid ranges for removal and insertion.
	if minRemovals >= maxRemovals || minDifference >= maxDifference {
		panic("max must be greater than min for both removal and insertion ranges")
	}

	// This variable will track the cycle number.
	cycleNumber := 0

	// This variable will keep track of the total operation count.
	var currentIncrement int64 = 0

	// Continue adding insertion and deletion patterns until reaching the target operation count.
	for currentIncrement < int64(randomTotalCount) {
		// Generate random values within the specified ranges for removals and insertions.
		removals := minRemovals + rand.Int63n(maxRemovals-minRemovals)
		difference := minDifference + rand.Int63n(maxDifference-minDifference)

		// Add a test stage with the generated insertion and deletion counts.
		testStages = append(testStages, stage{
			StageSummary: "Stage " + strconv.Itoa(cycleNumber),
			Op: struct {
				InsertAction int64
				DeleteAction int64
			}{
				InsertAction: removals + difference,
				DeleteAction: removals,
			},
			Repeat: 1,
		})

		// Increment the total count with the number of insertions performed in this cycle.
		currentIncrement += difference
		cycleNumber++
	}

	return testStages
}
