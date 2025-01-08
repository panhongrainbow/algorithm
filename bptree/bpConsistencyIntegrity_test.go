package bpTree

import (
	"fmt"
	"github.com/panhongrainbow/algorithm/randhub"
	"github.com/panhongrainbow/algorithm/testplan"
	bptestModel1 "github.com/panhongrainbow/algorithm/testplan/bptestplan/model1"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// =====================================================================================================================
//                  ⚗️ Consistency Integrity Test (B Plus Tree)
// =====================================================================================================================
// 🧪 The B Plus Tree unit test is designed to validate the tree’s consistency
// and integrity through bulk data insertion and deletion.
// 🧪 The test begins by inserting a large volume of data into the tree,
// followed by a complete deletion of all data, checking if the tree
// returns to its initial empty state to verify correctness.
// 🧪 Indexing errors in the B Plus Tree can lead to serious issues, such as
// being unable to find specific data or failing to delete data properly.
// 🧪 The test ensures the accuracy of indexing to prevent inconsistencies
// that might result in data operation failures.

// ⚗️ This code defines three constants used for generating random numbers in a test.
const (
	// 🧪 randomTotalCount represents the number of elements to be generated for random testing.
	randomTotalCount int64 = 7500000 // 20000000 // 147169280 // 7500000

	// 🧪 randomMin represents the minimum value for generating random numbers.
	randomMin int64 = 10

	// 🧪 random number hit collision percentage.
	randomHitCollisionPercentage int64 = 70

	// 🧪 randomMax represents the maximum value for generating random numbers.
	randomMax = randomTotalCount/randomHitCollisionPercentage*100 + randomMin
)

// Test_Check_BpTree_ConsistencyIntegrity 🧫 validates consistency and integrity by inserting and then deleting large data volumes
// to check if the tree returns to an empty state, ensuring indexing accuracy to prevent data operation failures.
func Test_Check_BpTree_ConsistencyIntegrity(t *testing.T) {
	testMode0Name := "Mode 0: Testing"
	t.Run(testMode0Name, func(t *testing.T) {
		model1 := &bptestModel1.BpTestModel1{RandomTotalCount: uint64(randomTotalCount)}

		dataSet, err := model1.GenerateRandomSet(1, 11)
		assert.NoError(t, err)

		fmt.Println(dataSet[0])
	})
	testMode1Name := "Mode 1: Bulk Insert/Delete"
	t.Run(testMode1Name, func(t *testing.T) {
		// Test case for bulk insert and delete operations on the B Plus tree.

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
		for i := 0; i < 5; i++ {
			// Create a test plan for bulk insert and delete operations.
			choosePlan := testplan.BpTreeProcess{
				RandomTotalCount: randomTotalCount, // Number of elements to generate for random testing.
			}
			testPlan := choosePlan.PlanMaxInsertDelete()

			// Generate a list of unique numbers for bulk insertion.
			bulkAdd, err := randhub.GenerateUniqueNumbers(uint64(randomTotalCount), randomMin, randomMax)
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
				uint32(randomTotalCount*2),                                        // Total number of operations.
				70,                                                                // Progress bar width.
				utilhub.WithTracking(5),                                           // Update interval.
				utilhub.WithTimeZone("Asia/Taipei"),                               // Time zone.
				utilhub.WithTimeControl(500),                                      // Update interval in milliseconds.
				utilhub.WithDisplay(utilhub.BrightCyan),                           // Display style.
			)

			// Start the progress bar printer in a separate goroutine.
			go func() {
				progressBar.ListenPrinter()
			}()

			// Initialize a new B Plus tree with a specified order.
			root := NewBpTree(bpTreeWidth + i)

			// Perform bulk insertion of generated numbers.
			for j := 0; j < int(randomTotalCount); j++ {
				// Insert a new value into the B Plus tree.
				root.InsertValue(BpItem{Key: bulkAdd[j]})
				// Update the progress bar.
				progressBar.UpdateBar()
			}

			// Perform bulk deletion of shuffled numbers.
			for k := 0; k < int(randomTotalCount); k++ {
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
				RandomTotalCount: randomTotalCount, // Number of elements to generate for random testing.
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
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

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
			_, removeAll := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, 0, 0, true)
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
				RandomTotalCount: randomTotalCount, // Number of elements to generate for random testing.
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
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

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
			_, removeAll := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, 0, 0, true)
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
				RandomTotalCount: randomTotalCount, // Number of elements to generate for random testing.
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
				batchInsert, batchRemove := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, int(testPlan[j].ChangePattern[0]), -1*int(testPlan[j].ChangePattern[1]), false)

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
			_, removeAll := pool.GenerateUniqueInt64Numbers(randomMin, randomMax, 0, 0, true)
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

	// Automated random testing for B+ tree.
	/*t.Run("Manually Identify B Plus Tree Operation Errors", func(t *testing.T) {
		// 数量二百的例子
		// Generate random data for insertion.
		var randomNumbers = []int64{1538, 1064, 249, 1966, 778, 1046, 1764, 1797, 847, 726, 1212, 119, 1063, 266, 1622, 511, 1609, 450, 1011, 707, 1425, 1045, 821, 1294, 1154, 1723, 1349, 1499, 230, 1320, 312, 917, 845, 1738, 1462, 236, 320, 1381, 409, 1805, 1709, 943, 1879, 69, 211, 367, 898, 1700, 1234, 395, 710, 1196, 1526, 384, 509, 1962, 1456, 205, 1830, 576, 587, 419, 1252, 1091, 346, 1066, 1876, 1088, 351, 1031, 1568, 1233, 761, 715, 691, 1368, 1739, 1314, 1197, 224, 1049, 1060, 1036, 1420, 567, 1305, 1618, 1557, 919, 115, 155, 1601, 874, 540, 260, 892, 1423, 794, 362, 484, 868, 1945, 958, 969, 1977, 905, 229, 1914, 376, 736, 156, 1105, 530, 1629, 405, 1398, 1706, 1691, 1683, 743, 1971, 279, 1256, 247, 1745, 785, 1119, 1513, 1078, 879, 1556, 1804, 1873, 388, 1418, 1880, 1362, 1392, 611, 930, 1240, 1571, 502, 1013, 439, 1581, 1881, 114, 235, 1703, 1341, 591, 393, 488, 538, 1925, 624, 1975, 1536, 1654, 89, 1902, 495, 1944, 674, 942, 1248, 1250, 660, 1928, 527, 1017, 1161, 1682, 71, 1807, 158, 1768, 435, 1623, 458, 1630, 125, 151, 67, 1087, 1820, 1009, 1506, 1823, 494, 863, 1439, 887, 765, 1600, 1428, 671, 1608, 1394, 583, 1288, 1824, 1737, 1180, 416, 1350, 1867, 714, 687, 138, 1535, 701, 614, 411, 1694, 522, 1913, 1328, 23, 1767, 87, 1365, 1441, 314, 1868, 262, 857, 54, 1055, 1765, 282, 1681, 43, 1325, 363, 1577, 396, 1302, 1023, 1190, 760, 417, 276, 1194, 1489, 1220, 1806, 1487, 275, 1659, 1859, 1777, 1163, 1204, 575, 1175, 947, 870, 163, 60, 104, 753, 1217, 748, 1244, 1758, 1658, 440, 1071, 888, 1592, 1040, 639, 601, 1417, 222, 118, 1862, 73, 605, 1003, 1177, 1152, 291, 1665, 1126, 330, 1464, 830, 762, 677, 933, 1301, 473, 1385, 652, 193, 1552, 621, 1888, 323, 1293, 1626, 597, 818, 343, 307, 940, 574, 446, 584, 1891, 1249, 97, 1690, 820, 918, 1313, 521, 935, 1183, 1229, 1253, 474, 1108, 1831, 1840, 582, 829, 1616, 34, 1363, 498, 1095, 675, 953, 1632, 1263, 103, 1653, 669, 453, 1751, 1495, 1171, 1704, 113, 1118, 1168}
		// Generate random data for deletion.
		var shuffledNumbers = []int64{236, 1571, 671, 1880, 669, 674, 715, 396, 119, 1806, 1758, 388, 363, 1420, 701, 193, 1626, 1777, 863, 691, 1709, 320, 1305, 857, 1196, 1556, 905, 1313, 458, 1234, 1428, 1204, 346, 1071, 151, 249, 918, 947, 1623, 343, 118, 1608, 97, 114, 887, 1618, 1629, 892, 1040, 743, 1171, 1745, 1423, 494, 113, 1925, 1694, 488, 276, 1737, 1658, 1064, 1609, 1240, 1363, 820, 1055, 1977, 611, 1046, 1060, 942, 760, 230, 1881, 1487, 247, 23, 1820, 1600, 1706, 540, 1088, 1738, 1513, 639, 235, 1031, 1190, 1108, 1249, 1263, 1723, 1876, 1462, 222, 1581, 1003, 1557, 1879, 1293, 1368, 1456, 1830, 1392, 34, 1506, 1180, 821, 266, 1441, 314, 940, 587, 1017, 1683, 1212, 1301, 829, 1690, 830, 1764, 1536, 1653, 1183, 1256, 1862, 1091, 262, 502, 1768, 1700, 1341, 1294, 384, 614, 1868, 1013, 205, 419, 1385, 943, 1233, 1804, 435, 1823, 60, 710, 1681, 687, 930, 376, 707, 511, 498, 919, 601, 71, 393, 440, 762, 1418, 958, 1632, 1797, 1320, 1654, 1049, 509, 1499, 591, 1288, 1244, 1194, 446, 1867, 933, 1126, 1105, 584, 367, 1568, 395, 583, 69, 411, 1577, 1161, 473, 1622, 484, 453, 1439, 1066, 1751, 1045, 43, 785, 953, 1592, 530, 1119, 156, 163, 1250, 158, 1944, 291, 527, 621, 879, 409, 870, 1767, 474, 417, 522, 73, 1220, 1425, 115, 1394, 1163, 1036, 736, 1914, 675, 888, 1975, 495, 1362, 362, 1859, 761, 576, 935, 868, 1704, 1739, 1349, 1023, 1552, 969, 1526, 1253, 597, 67, 1913, 765, 351, 1302, 575, 1177, 1325, 103, 605, 1118, 1252, 224, 818, 794, 1962, 1009, 1011, 138, 1945, 1175, 1971, 1152, 1495, 574, 1682, 714, 748, 1807, 624, 521, 677, 1902, 260, 1538, 1398, 1154, 229, 1840, 1616, 1873, 1805, 282, 778, 1229, 567, 652, 1765, 323, 450, 898, 1966, 1703, 1381, 1691, 1888, 874, 1314, 87, 753, 845, 1095, 1328, 439, 330, 1350, 1168, 1197, 1659, 726, 307, 405, 104, 125, 1489, 1217, 279, 1824, 1365, 660, 1601, 1417, 1063, 582, 155, 89, 1665, 1087, 917, 275, 1248, 847, 1891, 211, 1464, 1630, 312, 1535, 1928, 54, 416, 1078, 538, 1831}

		// Initialize B plus tree.
		root := NewBpTree(3)

		// Start inserting data.
		for i := 0; i < randomQuantity; i++ {
			// Insert data entries continuously.
			root.InsertValue(BpItem{Key: randomNumbers[i]})
		}

		// Start deleting data.
		for i := 0; i < randomQuantity; i++ {

			// 中断检验
			value := shuffledNumbers[i]
			fmt.Println(i, value)
			if shuffledNumbers[i] == 1824 {
				fmt.Println(">>>>> !")
				fmt.Print()
			}

			/*if shuffledNumbers[i] == 1365 {
				root.root.Index[1] = 1078
			}

			deleted, _, _, err := root.RemoveValue(BpItem{Key: shuffledNumbers[i]})

			if deleted == false {
				fmt.Println("Breakpoint: Data deletion not successful. 💢 The number is ", shuffledNumbers[i], i)
			}
			if err != nil {
				fmt.Println("Breakpoint: Deletion encountered an error. 💢 The number is ", shuffledNumbers[i], i)
			}
		}

		fmt.Println()
	})*/
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

/*func Test_Check_inode(t *testing.T) {
	root := NewBpTree(5)

	// Set up a top-level index node.
	root.root = &BpIndex{
		Index: []int64{288794, 339101, 460280},
		DataNodes: []*BpData{{
			Previous: nil,
			Next:     nil,
			Items:    []BpItem{{Key: 46911}, {Key: 54204}},
		}, {
			Previous: nil,
			Next:     nil,
			Items:    []BpItem{{Key: 288794}},
		}, {
			Previous: nil,
			Next:     nil,
			Items:    []BpItem{{Key: 339101}, {Key: 375797}},
		}, {
			Previous: nil,
			Next:     nil,
			Items:    []BpItem{{Key: 460280}, {Key: 468483}, {Key: 480770}},
		},
		},
	}

	root.InsertValue(BpItem{Key: 234438})
	root.InsertValue(BpItem{Key: 419488})
	root.InsertValue(BpItem{Key: 331451})

	root.root.Print()
}*/

/*func Test_Check_Continuity(t *testing.T) {
	root := NewBpTree(3)

	head := root.root.BpDataHead()

	root.InsertValue(BpItem{Key: 46911})
	root.InsertValue(BpItem{Key: 54204})
	root.InsertValue(BpItem{Key: 288794})
	root.InsertValue(BpItem{Key: 339101})
	root.InsertValue(BpItem{Key: 375797})
	root.InsertValue(BpItem{Key: 460280})
	root.InsertValue(BpItem{Key: 468483})
	root.InsertValue(BpItem{Key: 480770})

	// head.Items[0].Key = 1
	head.Next.Next.Next.Items[0].Key = 10

	root.CheckAndSwapRightContinuity()

	root.root.Print()
}*/

/*func printGoroutines() {
	buf := make([]byte, 1<<16) // 分配足够大的 buffer
	stackSize := runtime.Stack(buf, true)
	fmt.Printf("=== Goroutines Info ===\n%s\n", buf[:stackSize])
}*/
