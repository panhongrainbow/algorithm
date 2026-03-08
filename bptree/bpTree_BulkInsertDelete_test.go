package bpTree

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	bptestBulkInsertDelete "github.com/panhongrainbow/go-algorithm/testdata/bulkInsertDelete"
	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy BulkInsertDelete (Bulk Operation Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// prepare_BulkInsertDelete : prepares test data for BulkInsertDelete.
// verify_BulkInsertDelete : validates the test data.
// run_BulkInsertDelete : runs the test cases.
// =====================================================================================================================

// prepareBulkInsertDelete 🧫 prepares test data for BulkInsertDelete.
func prepareBulkInsertDelete(t *testing.T) {

	// === Init test model and record file ===

	// Create model 1 with specified data count.
	bptest1 := &bptestBulkInsertDelete.BpTestBulkInsertDelete{}

	// Create an empty record file.
	err := recordDir.Touch("BulkInsertDelete.do_not_open")
	require.NoError(t, err, "failed to create record file")

	// === Generate test data ===

	// Generate a random set: half positive, half negative.
	testDataSet, err := bptest1.GenerateRandomSet(uint64(unitTestConfig.Parameters.RandomMin), uint64(unitTestConfig.Parameters.RandomHitCollisionPercentage))
	require.NoError(t, err, "failed to generate test data")

	// === Set write parameters ===

	const (
		spliceBlockLength = 300
		spliceBlockWidth  = 100
	)

	// === Write data with Linux splice stream ===

	err = recordDir.LinuxSpliceProgressStreamWrite(
		testDataSet,
		"BulkInsertDelete.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"Bulk InsertDelete - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}

// verifyBulkInsertDelete 🧫 checks the test data set for BulkInsertDelete.
func verifyBulkInsertDelete(t *testing.T) {
	// Read test data with progress bar.
	testDataSet, err := recordDir.ReadAllBytesWithProgress(
		uint32(unitTestConfig.Parameters.RandomTotalCount),
		"BulkInsertDelete.do_not_open", 800,
		binary.LittleEndian,
		"Bulk InsertDelete - read test data",
		utilhub.BrightCyan,
		70,
	)

	// Init test model.
	bptest1 := &bptestBulkInsertDelete.BpTestBulkInsertDelete{}

	// Validate test data.
	err = bptest1.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to validate test data")
}

// runBulkInsertDelete 🧫 runs the actual test cases for BulkInsertDelete.
func runBulkInsertDelete(t *testing.T) {
	for bpWidth := 0; bpWidth < len(unitTestConfig.Parameters.BpWidth); bpWidth++ {
		_runBulkInsertDelete(t, bpWidth)
	}
}

// _runBulkInsertDelete 🧫 runs the actual test cases for BulkInsertDelete.
func _runBulkInsertDelete(t *testing.T, bpWidth int) {
	dtatChan, errChan, finsishChan := recordDir.ReadBytesInChunksWithProgress("BulkInsertDelete.do_not_open", 8, binary.LittleEndian)

	root := NewBpTree(unitTestConfig.Parameters.BpWidth[bpWidth])

	// testBulkInsertDeleteName := "Execution; Width: " + strconv.Itoa(unitTestConfig.Parameters.BpWidth[bpWidth])
	testBulkInsertDeleteName := fmt.Sprintf("Bulk Insert/Delete - run; Width: %3d", unitTestConfig.Parameters.BpWidth[bpWidth])

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		testBulkInsertDeleteName,
		// "BulkInsertDelete: Execution   ",                             // Progress bar title.
		uint32(unitTestConfig.Parameters.RandomTotalCount), // Total number of operations.
		70,                                       // Progress bar width.
		utilhub.WithTracking(5),                  // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),      // Time zone.
		utilhub.WithTimeControl(500),             // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightGreen), // Display style.
	)

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

Loop:
	for {
		select {
		case data := <-dtatChan:
			for j := 0; j < len(data); j++ {
				if data[j] >= 0 {
					root.InsertValue(BpItem{Key: data[j]})
					progressBar.UpdateBar()
				}
				if data[j] < 0 {
					deleted, _, _, err := root.RemoveValue(BpItem{Key: -1 * data[j]})
					require.True(t, deleted)
					require.NoError(t, err)
					progressBar.UpdateBar()
				}
			}
		case err := <-errChan:
			fmt.Println(err)
		case <-finsishChan:
			break Loop
		}
	}

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Print a final report.
	err := progressBar.Report(len(testBulkInsertDeleteName + "; Width: XX"))
	assert.NoError(t, err)

	// Print the B Plus tree structure.
	root.root.Print()
}
