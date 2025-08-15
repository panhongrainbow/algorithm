package model2

import (
	"math/rand"
	"strconv"

	"github.com/panhongrainbow/algorithm/utilhub"
)

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
	// UseFixedData bool

	// DataSource contains the data used in this stage, which can include
	// either predefined values or randomly generated data.
	// DataSource []int64

	// IsFinalStage indicates whether this is the final stage of the test.
	// If true, no additional stages will follow.
	// IsFinalStage bool
}

// RandomizedBoundary 🧮 plans repeated insertion and deletion of data in the boundary of B Plus tree.
func (model2 *BpTestModel2) RandomizedBoundary(minRemovals, maxRemovals, minDifference, maxDifference int64) (testStages []EachBpTestStage) {
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
		testStages = append(testStages, EachBpTestStage{
			Description:    "Cycle " + strconv.Itoa(cycleNumber),
			ChangePattern:  []int64{removals + difference, -1 * removals},
			ExecutionCycle: 1,
			// UseFixedData:   false,
		})

		// Increment the total count with the number of insertions performed in this cycle.
		currentIncrement += difference
		cycleNumber++
	}

	return testStages
}

func (model2 *BpTestModel2) TotalOperation(testStages []EachBpTestStage) int64 {
	var totalOperation int64
	for i := 0; i < len(testStages); i++ {

		// Determine the number of times to repeat the test execution cycle.
		// If the ExecutionCycle is greater than 1, use that value; otherwise, default to 1.
		repeatTime := 1
		if testStages[i].ExecutionCycle > 1 {
			// Override the default repeat time with the specified ExecutionCycle value.
			repeatTime = testStages[i].ExecutionCycle
		}

		totalOperation += testStages[i].ChangePattern[0] * int64(repeatTime) * 2
	}
	return totalOperation
}
