package bpTree

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	bptestModel1 "github.com/panhongrainbow/algorithm/testdata/model1"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy Mode 1 (Test Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// [prepare_Mode1] prepares test data for Mode 1.
// [verify_Mode1] validates the test data.
// [run_Mode1] runs the test cases.
// =====================================================================================================================

// prepareMode1 🧫 prepares test data for Mode 1.
func prepareMode1(t *testing.T) {

	// === Init test model and record file ===

	// Create model 1 with specified data count.
	bptest1 := &bptestModel1.BpTestModel1{}

	// Create an empty record file.
	err := recordDir.Touch("mode1.do_not_open")
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
		"mode1.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"Mode 1: Bulk Insert/Delete - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}

// verifyMode1 🧫 checks the test data set for Mode 1.
func verifyMode1(t *testing.T) {
	// Read test data with progress bar.
	testDataSet, err := recordDir.ReadAllBytesWithProgress(
		uint32(unitTestConfig.Parameters.RandomTotalCount),
		"mode1.do_not_open", 800,
		binary.LittleEndian,
		"Mode 1: Bulk Insert/Delete - read test data",
		utilhub.BrightCyan,
		70,
	)

	// Init test model.
	bptest1 := &bptestModel1.BpTestModel1{}

	// Validate test data.
	err = bptest1.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to validate test data")
}

// runMode1 🧫 runs the actual test cases for Mode 1.
func runMode1(t *testing.T) {
	for bpWidth := 0; bpWidth < len(unitTestConfig.Parameters.BpWidth); bpWidth++ {
		_runMode1(t, bpWidth)
	}
}

// _runMode1 🧫 runs the actual test cases for Mode 1.
func _runMode1(t *testing.T, bpWidth int) {
	dtatChan, errChan, finsishChan := recordDir.ReadBytesInChunksWithProgress("mode1.do_not_open", 8, binary.LittleEndian)

	root := NewBpTree(unitTestConfig.Parameters.BpWidth[bpWidth])

	// testMode1Name := "Mode 1: Execution; Width: " + strconv.Itoa(unitTestConfig.Parameters.BpWidth[bpWidth])
	testMode1Name := fmt.Sprintf("Mode 1: Bulk Insert/Delete - run; Width: %3d", unitTestConfig.Parameters.BpWidth[bpWidth])

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		testMode1Name,
		// "Mode 1: Execution   ",                             // Progress bar title.
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
	err := progressBar.Report(len(testMode1Name + "; Width: XX"))
	assert.NoError(t, err)

	// Print the B Plus tree structure.
	root.root.Print()
}
