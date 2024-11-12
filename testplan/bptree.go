package testplan

import (
	"math/rand"
	"strconv"
)

// =====================================================================================================================
//                   🧮 BpTree Algorithm Test Plan
// =====================================================================================================================
// ✏️ These functions are designed to systematically plan the data generation,
// batch insertion, and deletion processes for algorithm testing.
// ✏️ They ensure the test process is structured, clear, and easy to maintain,
// while effectively covering various test scenarios.
// ✏️ By organizing the workflow in advance, this approach enhances the
// test's overall effectiveness and reliability.

// EachBpTestStage Slice 🧮 will encompass a wide range of combinations, including boundary condition tests, consistency tests, balance tests, and more.
type EachBpTestStage struct {
	// Description provides a summary or purpose of the stage, explaining its testing goal.
	Description string

	// ChangePattern records the specific data changes in this stage,
	// such as inserting 20 items and then deleting 5 items.
	ChangePattern []int64

	// ExecutionCycle defines the number of times this stage will repeat.
	ExecutionCycle int

	// UseFixedData specifies whether fixed data should be used.
	// If true, DataSource will be used for this stage.
	UseFixedData bool

	// DataSource contains the data used in this stage, which can include
	// either predefined values or randomly generated data.
	DataSource []int64

	// IsFinalStage indicates whether this is the final stage of the test.
	// If true, no additional stages will follow.
	IsFinalStage bool

	// Statistic holds various metrics and descriptions for this stage.
	Statistic EachTestStageStatistic
}

// EachTestStageStatistic 🧮 represents statistical data collected for each test stage.
type EachTestStageStatistic struct {
	// MaxDataAmount represents the maximum data amount to reach in this stage.
	MaxDataAmount int64

	// MinDataAmount represents the minimum data amount to reach in this stage.
	MinDataAmount int64

	// DataAmountRange is the range of data amount, calculated as MaxDataAmount minus MinDataAmount.
	DataAmountRange int64
}

/*
	🧮 BpTree Test Process Calculator
	The process calculator is a simple tool designed to help users plan and calculate workflows.
	It breaks down complex operations into a series of simple steps and calculates the operations or results required for each step.
*/

// BpTreeProcess 🧮 stores common parameters for the process calculator.
type BpTreeProcess struct {
	// randomTotalCount represents the number of elements to be generated for random testing.
	RandomTotalCount int64
}

// Mode 1: Bulk Insert/Delete

// PlanMaxInsertDelete 🧮 generates a plan for sudden bulk insertion followed by bulk deletion of data in the B Plus tree.
func (bPlan BpTreeProcess) PlanMaxInsertDelete() []EachBpTestStage {
	return []EachBpTestStage{
		{
			Description:    "Session 1 : Bulk Insert",
			ChangePattern:  []int64{bPlan.RandomTotalCount},
			ExecutionCycle: 1,
			UseFixedData:   false,
		},
		{
			Description:    "Session 2 : Bulk Delete",
			ChangePattern:  []int64{-1 * bPlan.RandomTotalCount},
			ExecutionCycle: 1,
			UseFixedData:   false,
		},
	}
}

// InsertionDeletionLoad 🧮 plans repeated insertion and deletion of data in the B Plus tree.
func (bPlan BpTreeProcess) InsertionDeletionLoad(minRemovals, maxRemovals, minDifference, maxDifference int64) (testStages []EachBpTestStage) {
	// Perform boundary check outside of the loop to ensure valid ranges for removal and insertion
	if minRemovals >= maxRemovals || minDifference >= maxDifference {
		panic("max must be greater than min for both removal and insertion ranges")
	}

	cycleNumber := 0
	var currentIncrement int64 = 0

	// Continue adding insertion and deletion patterns until reaching the target operation count
	for currentIncrement < bPlan.RandomTotalCount {
		// Generate random values within the specified ranges for removals and insertions
		removals := minRemovals + rand.Int63n(maxRemovals-minRemovals)
		difference := minDifference + rand.Int63n(maxDifference-minDifference)

		// Add a test stage with the generated insertion and deletion counts
		testStages = append(testStages, EachBpTestStage{
			Description:    "Cycle " + strconv.Itoa(cycleNumber),
			ChangePattern:  []int64{removals + difference, -1 * removals},
			ExecutionCycle: 1,
			UseFixedData:   false,
		})

		// Increment the total count with the number of insertions performed in this cycle
		currentIncrement += difference
		cycleNumber++
	}

	return testStages
}
func (bPlan BpTreeProcess) TotalOperation(testStages []EachBpTestStage) int64 {
	var totalOperation int64
	for i := 0; i < len(testStages); i++ {
		totalOperation += testStages[i].ChangePattern[0] * 2
	}
	return totalOperation
}
