package model2

import (
	"math/rand"
	"strconv"

	"github.com/panhongrainbow/algorithm/utilhub"
)

// StagePlan defines a test stage with operations and repeat count.
type StagePlan struct {
	// StageSummary provides a short description of this test stage.
	StageSummary string

	// OperationPlan specifies the sequence of operations (e.g., insert/delete counts).
	OperationPlan []int64

	// RepeatCount indicates how many times this stage will be executed.
	RepeatCount int
}

// CountOps returns the total number of insert/delete operations across all stages.
//
// The sum of all OperationPlans across stages is defined to be zero (total inserts equal total deletes).
// Therefore, the total number of operations can be calculated as:
//
//	insertions * 2 * RepeatCount
//
// where OperationPlan[0] is used as the insertion count.
func (model2 *BpTestModel2) CountOps(stages []StagePlan) int64 {
	var totalOps int64
	for _, stage := range stages {
		if stage.RepeatCount > 1 {
			totalOps += stage.OperationPlan[0] * int64(stage.RepeatCount) * 2
		}
	}
	return totalOps
}

// RandomizedBoundary 🧮 plans repeated insertion and deletion of data in the boundary of B Plus tree.
func (model2 *BpTestModel2) RandomizedBoundary(minRemovals, maxRemovals, minDifference, maxDifference int64) (testStages []StagePlan) {
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
		testStages = append(testStages, StagePlan{
			StageSummary:  "Stage " + strconv.Itoa(cycleNumber),
			OperationPlan: []int64{removals + difference, -1 * removals},
			RepeatCount:   1,
		})

		// Increment the total count with the number of insertions performed in this cycle.
		currentIncrement += difference
		cycleNumber++
	}

	return testStages
}
