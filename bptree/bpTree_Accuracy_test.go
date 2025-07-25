package bpTree

// =====================================================================================================================
//                  ⚗️ Consistency Integrity Test ( [B Plus Tree] ) - B加树 主要测试
// =====================================================================================================================
// 🧪 B Plus Tree unit test validates structure via bulk insert and delete.
// 🧪 Inserts large data, then deletes all to check if tree resets to empty.
// 🧪 Indexing errors may cause data loss or deletion failures.
// 🧪 Ensures indexing accuracy for reliable operations.

// To run the test, run the following command:
//
// cd /home/panhong/go/src/github.com/panhongrainbow/algorithm/bptree
// go clean -cache
// go test -v . -timeout=0 -run Test_Check_BpTree_Accuracies

// =====================================================================================================================

import (
	"fmt"
	"github.com/panhongrainbow/algorithm/randhub"
	"github.com/panhongrainbow/algorithm/testplan"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var (
	// 🧪 Create a config struct for B plus tree unit testing and parse default values.
	unitTestConfig = utilhub.BptreeUnitTestConfig{}
	configParseErr = utilhub.ParseDefault(&unitTestConfig)

	// 🧪 Navigate to the project directory for test record storage.
	ProjectDir = utilhub.FileNode{}.Goto(unitTestConfig.Record.TestRecordPath)

	// 🧪 Create a subdirectory named with the current date under the project.
	recordDir = ProjectDir.MkDir(time.Now().Format("2006-01-02"))
)

// Test_Check_BpTree_Accuracy 🧫 checks if the tree resets after bulk insert/delete, ensuring indexing correctness.
func Test_Check_BpTree_Accuracies(t *testing.T) {

	t.Run("Pre-test checks", func(t *testing.T) {
		// Config must parse without error.
		require.NoError(t, configParseErr)

		// Record path must not be empty.
		require.NotEqual(t, "", ProjectDir.Path(), "record path is empty; check path creation")

		// Record subdirectory must not be empty.
		require.NotEqual(t, "", recordDir.Path(), "record date path is empty; check path creation")
	})

	testMode0Name := "Mode 0: Bulk Insert/Delete"
	t.Run(testMode0Name, func(t *testing.T) {
		// Prepare test data for mode 1.
		prepareMode1(t)

		// Verify test data for mode 1.
		verifyMode1(t)

		// Execute accuracy test for mode 1.
		runMode1(t)
	})

	testMode1Name := "Mode 1: Bulk Insert/Delete"
	t.Run(testMode1Name, func(t *testing.T) {
		// Test case for bulk insert and delete operations on the B Plus tree.

		// Initialize a random number generator.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Set the initial width of the B Plus tree.
		// The width of the B Plus Tree determines the number of keys that can be stored in each node.

		// The random width is between 3 and 12.
		// This is done to ensure that the number of keys in each node is varied,
		// which helps to check for errors in indexing.
		// The test is repeated five times, each time with an incremented width.
		// This includes testing with both odd and even widths.
		bpTreeWidth := rand.Intn(10) + 3

		// Perform tests with varying B Plus tree widths to ensure robustness.
		for i := 0; i < 5; i++ {
			// Create a test plan for bulk insert and delete operations.
			choosePlan := testplan.BpTreeProcess{
				RandomTotalCount: unitTestConfig.Parameters.RandomTotalCount, // Number of elements to generate for random testing.
			}
			testPlan := choosePlan.PlanMaxInsertDelete()

			// Generate a list of unique numbers for bulk insertion.
			bulkAdd, err := randhub.GenerateUniqueNumbers(uint64(unitTestConfig.Parameters.RandomTotalCount), unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax)
			if err != nil {
				// Panic if an error occurs during number generation.
				panic(err)
			}

			// Create a copy of the bulk insertion list and shuffle it for deletion.
			bulkDel := make([]int64, testPlan[0].ChangePattern[0])
			copy(bulkDel, bulkAdd)
			shuffleSlice(bulkDel, random)

			// Create a progress bar with optional configurations.
			progressBar, _ := utilhub.NewProgressBar(
				"Mode 1: Bulk Insert/Delete; Width: "+strconv.Itoa(bpTreeWidth+i), // Progress bar title.
				uint32(unitTestConfig.Parameters.RandomTotalCount*2),              // Total number of operations.
				70,                                      // Progress bar width.
				utilhub.WithTracking(5),                 // Update interval.
				utilhub.WithTimeZone("Asia/Taipei"),     // Time zone.
				utilhub.WithTimeControl(500),            // Update interval in milliseconds.
				utilhub.WithDisplay(utilhub.BrightCyan), // Display style.
			)

			// Start the progress bar printer in a separate goroutine.
			go func() {
				progressBar.ListenPrinter()
			}()

			// Initialize a new B Plus tree with a specified order.
			root := NewBpTree(bpTreeWidth + i)

			// Perform bulk insertion of generated numbers.
			for j := 0; j < int(unitTestConfig.Parameters.RandomTotalCount); j++ {
				// Insert a new value into the B Plus tree.
				root.InsertValue(BpItem{Key: bulkAdd[j]})
				// Update the progress bar.
				progressBar.UpdateBar()
			}

			// Perform bulk deletion of shuffled numbers.
			for k := 0; k < int(unitTestConfig.Parameters.RandomTotalCount); k++ {
				// Remove a value from the B Plus tree.
				deleted, _, _, err := root.RemoveValue(BpItem{Key: bulkDel[k]})
				// Update the progress bar.
				progressBar.UpdateBar()

				// Check for errors during deletion.
				if err != nil {
					// Panic with detailed error message about the failure during deletion.
					panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", bulkDel[k], k, err))
				}

				// Check if deletion was successful.
				if deleted == false {
					// Panic with detailed error message indicating deletion was not successful.
					panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", bulkDel[k], k))
				}
			}

			// Mark the progress bar as complete.
			progressBar.Complete()

			// Wait for the progress bar printer to stop.
			<-progressBar.WaitForPrinterStop()

			// Print a final report.
			err = progressBar.Report(len(testMode1Name + "; Width: XX"))
			assert.NoError(t, err)

			// Print the B Plus tree structure.
			root.root.Print()
		}
	})
	testMode2Name := "Mode 2: Randomized Boundary Test"
	t.Run(testMode2Name, func(t *testing.T) {
		// By repeatedly performing insert and delete operations, we can assess the system's
		// stability, performance, correctness, and handling of edge cases when dealing with a dynamic dataset.

		// Initialize a random number generator.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Set the initial width of the B Plus tree.
		// The width of the B Plus tree determines the number of keys that can be stored in each node.

		// The random width is between 3 and 12.
		// This is done to ensure that the number of keys in each node is varied,
		// which helps to check for errors in indexing.
		// The test is repeated five times, each time with an incremented width.
		// This includes testing with both odd and even widths.
		bpTreeWidth := rand.Intn(10) + 3

		// Perform tests with varying B Plus tree widths to ensure robustness.
		// We need at least 2 iterations to cover both odd and even BpTree widths.
		for i := 0; i < 5; i++ {
			// Create a test plan for bulk insert and delete operations.
			choosePlan := testplan.BpTreeProcess{
				RandomTotalCount: unitTestConfig.Parameters.RandomTotalCount, // Number of elements to generate for random testing.
			}
			testPlan := choosePlan.RandomizedBoundary(5, 50, 10, 20)

			// Create a progress bar with optional configurations.
			progressBar, _ := utilhub.NewProgressBar(
				"Mode 2: Randomized Boundary Test; Width: "+strconv.Itoa(bpTreeWidth+i), // Progress bar title.
				uint32(choosePlan.TotalOperation(testPlan)),                             // Total number of operations.
				70,                                      // Progress bar width.
				utilhub.WithTracking(5),                 // Update interval.
				utilhub.WithTimeZone("Asia/Taipei"),     // Time zone.
				utilhub.WithTimeControl(500),            // Update interval in milliseconds.
				utilhub.WithDisplay(utilhub.BrightCyan), // Display style.
			)

			// Start the progress bar printer in a separate goroutine.
			go func() {
				progressBar.ListenPrinter()
			}()

			// Initialize a test pool for generating random numbers.
			pool := randhub.NewDoublePool()

			// Initialize a new B Plus tree with a specified order.
			root := NewBpTree(bpTreeWidth + i)

			// Iterate through the test plan for bulk insert and delete operations in order to test stability and consistency.
			for j := 0; j < len(testPlan); j++ {
				// Generate random numbers for bulk insertion and deletion.
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

				// Create a copy of the bulk insertion list and shuffle it for deletion.
				shuffleSlice(batchInsert, random)
				shuffleSlice(batchRemove, random)

				// Perform bulk insertion of generated numbers.
				for k := 0; k < int(testPlan[j].ChangePattern[0]); k++ {
					// Insert a new value into the B Plus tree.
					root.InsertValue(BpItem{Key: batchInsert[k]})
					// Update the progress bar.
					progressBar.UpdateBar()
				}

				// Perform bulk deletion of shuffled numbers.
				for l := 0; l < -1*int(testPlan[j].ChangePattern[1]); l++ {
					// Remove a value from the B Plus tree.
					deleted, _, _, err := root.RemoveValue(BpItem{Key: batchRemove[l]})
					// Update the progress bar.
					progressBar.UpdateBar()

					// Check for errors during deletion.
					if err != nil {
						// Panic with detailed error message about the failure during deletion.
						panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", batchRemove[l], l, err))
					}

					// Check if deletion was successful.
					if deleted == false {
						// Panic with detailed error message indicating deletion was not successful.
						panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", batchRemove[l], l))
					}
				}
			}

			// Delete all data from the B Plus tree.
			_, removeAll := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, 0, 0, true)
			for m := 0; m < len(removeAll); m++ {
				deleted, _, _, err := root.RemoveValue(BpItem{Key: removeAll[m]})
				progressBar.UpdateBar()
				if deleted == false {
					// Panic with detailed error message about the failure during deletion.
					panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", removeAll[m], m, err))
				}
				if err != nil {
					// Panic with detailed error message indicating deletion was not successful.
					panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", removeAll[m], m))
				}
			}

			// Mark the progress bar as complete.
			progressBar.Complete()

			// Wait for the progress bar printer to stop.
			<-progressBar.WaitForPrinterStop()

			// Print a final report.
			err := progressBar.Report(len(testMode2Name + "; Width: XX"))
			assert.NoError(t, err)

			// Print the B Plus tree structure.
			root.root.Print()
		}
	})
	testMode3Name := "Mode 3: Gradual Boundary Test"
	t.Run(testMode3Name, func(t *testing.T) {
		// By repeatedly performing insert and delete operations, we can assess the system's
		// stability, performance, correctness, and handling of edge cases when dealing with a dynamic dataset.

		// Initialize a random number generator.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Set the initial width of the B Plus tree.
		// The width of the B+ tree determines the number of keys that can be stored in each node.

		// The random width is between 3 and 12.
		// This is done to ensure that the number of keys in each node is varied,
		// which helps to check for errors in indexing.
		// The test is repeated five times, each time with an incremented width.
		// This includes testing with both odd and even widths.
		bpTreeWidth := rand.Intn(10) + 3

		// Perform tests with varying B Plus tree widths to ensure robustness.
		// We need at least 2 iterations to cover both odd and even BpTree widths.
		for i := 0; i < 5; i++ {
			// Create a test plan for bulk insert and delete operations.
			choosePlan := testplan.BpTreeProcess{
				RandomTotalCount: unitTestConfig.Parameters.RandomTotalCount, // Number of elements to generate for random testing.
			}
			testPlan := choosePlan.GradualBoundary(5, 50, 10, 20)

			// Create a progress bar with optional configurations.
			progressBar, _ := utilhub.NewProgressBar(
				"Mode 3: Gradual Boundary Test; Width: "+strconv.Itoa(bpTreeWidth+i), // Progress bar title.
				uint32(choosePlan.TotalOperation(testPlan)),                          // Total number of operations.
				70,                                      // Progress bar width.
				utilhub.WithTracking(5),                 // Update interval.
				utilhub.WithTimeZone("Asia/Taipei"),     // Time zone.
				utilhub.WithTimeControl(500),            // Update interval in milliseconds.
				utilhub.WithDisplay(utilhub.BrightCyan), // Display style.
			)

			// Start the progress bar printer in a separate goroutine.
			go func() {
				progressBar.ListenPrinter()
			}()

			// Initialize a test pool for generating random numbers.
			pool := randhub.NewDoublePool()

			// Initialize a new B Plus tree with a specified order.
			root := NewBpTree(bpTreeWidth + i)

			// Iterate through the test plan for bulk insert and delete operations in order to test stability and consistency.
			for j := 0; j < len(testPlan); j++ {
				// Generate random numbers for bulk insertion and deletion.
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

				// Create a copy of the bulk insertion list and shuffle it for deletion.
				shuffleSlice(batchInsert, random)
				shuffleSlice(batchRemove, random)

				// Perform bulk insertion of generated numbers.
				for k := 0; k < int(testPlan[j].ChangePattern[0]); k++ {
					// Insert a new value into the B Plus tree.
					root.InsertValue(BpItem{Key: batchInsert[k]})
					// Update the progress bar.
					progressBar.UpdateBar()
				}

				// Perform bulk deletion of shuffled numbers.
				for l := 0; l < -1*int(testPlan[j].ChangePattern[1]); l++ {
					// Remove a value from the B Plus tree.
					deleted, _, _, err := root.RemoveValue(BpItem{Key: batchRemove[l]})
					// Update the progress bar.
					progressBar.UpdateBar()

					// Check for errors during deletion.
					if err != nil {
						// Panic with detailed error message about the failure during deletion.
						panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", batchRemove[l], l, err))
					}

					// Check if deletion was successful.
					if deleted == false {
						// Panic with detailed error message indicating deletion was not successful.
						panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", batchRemove[l], l))
					}
				}
			}

			// Delete all data from the B Plus tree.
			_, removeAll := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, 0, 0, true)
			for m := 0; m < len(removeAll); m++ {
				deleted, _, _, err := root.RemoveValue(BpItem{Key: removeAll[m]})
				progressBar.UpdateBar()
				if deleted == false {
					// Panic with detailed error message about the failure during deletion.
					panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", removeAll[m], m, err))
				}
				if err != nil {
					// Panic with detailed error message indicating deletion was not successful.
					panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", removeAll[m], m))
				}
			}

			// Mark the progress bar as complete.
			progressBar.Complete()

			// Wait for the progress bar printer to stop.
			<-progressBar.WaitForPrinterStop()

			// Print a final report.
			err := progressBar.Report(len(testMode3Name + "; Width: XX"))
			assert.NoError(t, err)

			// Print the B Plus tree structure.
			root.root.Print()
		}
	})
	testMode4Name := "Mode 4: Redundant Operation"
	t.Run(testMode4Name, func(t *testing.T) {
		// By repeatedly performing insert and delete operations, we can assess the system's
		// stability, performance, correctness, and handling of edge cases when dealing with a dynamic dataset.

		// Initialize a random number generator.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Set the initial width of the B Plus tree.
		// The width of the B+ tree determines the number of keys that can be stored in each node.

		// The random width is between 3 and 12.
		// This is done to ensure that the number of keys in each node is varied,
		// which helps to check for errors in indexing.
		// The test is repeated five times, each time with an incremented width.
		// This includes testing with both odd and even widths.
		bpTreeWidth := rand.Intn(10) + 3

		// Perform tests with varying B Plus tree widths to ensure robustness.
		// We need at least 2 iterations to cover both odd and even BpTree widths.
		for i := 0; i < 5; i++ {
			// Create a test plan for bulk insert and delete operations.
			choosePlan := testplan.BpTreeProcess{
				RandomTotalCount: unitTestConfig.Parameters.RandomTotalCount, // Number of elements to generate for random testing.
			}
			testPlan := choosePlan.GradualBoundary(5, 50, 10, 20)

			// Create a progress bar with optional configurations.
			progressBar, _ := utilhub.NewProgressBar(
				"Mode 4: Gradual Boundary Test; Width: "+strconv.Itoa(bpTreeWidth+i), // Progress bar title.
				uint32(choosePlan.TotalOperation(testPlan)),                          // Total number of operations.
				70,                                      // Progress bar width.
				utilhub.WithTracking(5),                 // Update interval.
				utilhub.WithTimeZone("Asia/Taipei"),     // Time zone.
				utilhub.WithTimeControl(500),            // Update interval in milliseconds.
				utilhub.WithDisplay(utilhub.BrightCyan), // Display style.
			)

			// Start the progress bar printer in a separate goroutine.
			go func() {
				progressBar.ListenPrinter()
			}()

			// Initialize a test pool for generating random numbers.
			pool := randhub.NewDoublePool()

			// Initialize a new B Plus tree with a specified order.
			root := NewBpTree(bpTreeWidth + i)

			// Iterate through the test plan for bulk insert and delete operations in order to test stability and consistency.
			for j := 0; j < len(testPlan); j++ {
				// Generate random numbers for bulk insertion and deletion.
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

				// Create a copy of the bulk insertion list and shuffle it for deletion.
				shuffleSlice(batchInsert, random)
				shuffleSlice(batchRemove, random)

				// Perform bulk insertion of generated numbers.
				for k := 0; k < int(testPlan[j].ChangePattern[0]); k++ {
					// Insert a new value into the B Plus tree.
					root.InsertValue(BpItem{Key: batchInsert[k]})
					// Update the progress bar.
					progressBar.UpdateBar()
				}

				// Perform bulk deletion of shuffled numbers.
				for l := 0; l < -1*int(testPlan[j].ChangePattern[1]); l++ {
					// Remove a value from the B Plus tree.
					deleted, _, _, err := root.RemoveValue(BpItem{Key: batchRemove[l]})
					// Update the progress bar.
					progressBar.UpdateBar()

					// Check for errors during deletion.
					if err != nil {
						// Panic with detailed error message about the failure during deletion.
						panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", batchRemove[l], l, err))
					}

					// Check if deletion was successful.
					if deleted == false {
						// Panic with detailed error message indicating deletion was not successful.
						panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", batchRemove[l], l))
					}
				}

				for n := 0; n < 30; n++ {
					// Perform bulk insertion of generated numbers.
					for o := 0; o < -1*int(testPlan[j].ChangePattern[1]); o++ {
						// Insert a new value into the B Plus tree.
						root.InsertValue(BpItem{Key: batchRemove[o]})
						// Update the progress bar.
						// progressBar.UpdateBar()
					}

					// Perform bulk deletion of shuffled numbers.
					for p := 0; p < -1*int(testPlan[j].ChangePattern[1]); p++ {
						// Remove a value from the B Plus tree.
						deleted, _, _, err := root.RemoveValue(BpItem{Key: batchRemove[p]})
						// Update the progress bar.
						// progressBar.UpdateBar()

						// Check for errors during deletion.
						if err != nil {
							// Panic with detailed error message about the failure during deletion.
							panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", batchRemove[p], p, err))
						}

						// Check if deletion was successful.
						if deleted == false {
							// Panic with detailed error message indicating deletion was not successful.
							panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", batchRemove[p], p))
						}
					}
				}
			}

			// Delete all data from the B Plus tree.
			_, removeAll := pool.GenerateUniqueInt64Numbers(unitTestConfig.Parameters.RandomMin, unitTestConfig.Parameters.RandomMax, 0, 0, true)
			for m := 0; m < len(removeAll); m++ {
				deleted, _, _, err := root.RemoveValue(BpItem{Key: removeAll[m]})
				progressBar.UpdateBar()
				if deleted == false {
					// Panic with detailed error message about the failure during deletion.
					panic(fmt.Sprintf("Error during deletion: Failed to delete number %d at index %d. Error: %v", removeAll[m], m, err))
				}
				if err != nil {
					// Panic with detailed error message indicating deletion was not successful.
					panic(fmt.Sprintf("Error during deletion: Data deletion for number %d at index %d was not successful.", removeAll[m], m))
				}
			}

			// Mark the progress bar as complete.
			progressBar.Complete()

			// Wait for the progress bar printer to stop.
			<-progressBar.WaitForPrinterStop()

			// Print a final report.
			err := progressBar.Report(len(testMode4Name + "; Width: XX"))
			assert.NoError(t, err)

			// Print the B Plus tree structure.
			root.root.Print()
		}
	})
}

// shuffleSlice randomly shuffles the elements in the slice.
func shuffleSlice(slice []int64, rng *rand.Rand) {
	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := rng.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}
