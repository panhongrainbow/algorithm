package bpTree

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	bptestRandomizedBoundary "github.com/panhongrainbow/go-algorithm/testdata/randomizedBoundary"
	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy RandomizedBoundary (Boundary Test Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// prepare_RandomizedBoundary : prepares test data for RandomizedBoundary.
// verify_RandomizedBoundary : validates the test data.
// run_RandomizedBoundary : runs the test cases.
// =====================================================================================================================

// prepareRandomizedBoundary 🧫 prepares test data for RandomizedBoundary.
func prepareRandomizedBoundary(t *testing.T, recordDir utilhub.FileNode) {

	// === Init test model and record file ===

	// Create RandomizedBoundary test with specified data count.
	bptest := &bptestRandomizedBoundary.BpTestRandomizedBoundary{}

	// Create an empty record file.
	err := recordDir.Touch("RandomizedBoundary.do_not_open")
	require.NoError(t, err, "failed to create record file")

	// === Generate test data ===

	// Generate a random set: half positive, half negative.
	var testDataSet []int64
	testDataSet, err = bptest.GenerateRandomSet()
	require.NoError(t, err, "failed to generate test data")

	// === Set write parameters ===

	const (
		spliceBlockLength = 300
		spliceBlockWidth  = 100
	)

	// === Write data with Linux splice stream ===

	err = recordDir.LinuxSpliceProgressStreamWrite(
		testDataSet,
		"RandomizedBoundary.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"RandomizedBoundary - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}

// verifyRandomizedBoundary 🧫 checks the test data set for RandomizedBoundary.
func verifyRandomizedBoundary(t *testing.T, recordDir utilhub.FileNode, autoTestConfig utilhub.AutoConfigType) {
	// Read test data with progress bar.
	testDataSet, err := recordDir.ReadAllBytesWithProgress(
		uint32(autoTestConfig.Parameters.RandomTotalCount),
		"RandomizedBoundary.do_not_open", 800,
		binary.LittleEndian,
		"RandomizedBoundary - read test data",
		utilhub.BrightCyan,
		70,
	)

	// Init test model.
	bptest := &bptestRandomizedBoundary.BpTestRandomizedBoundary{}

	// Validate test data.
	err = bptest.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to validate test data")
}

// runRandomizedBoundary 🧫 runs the actual test cases for RandomizedBoundary.
func runRandomizedBoundary(t *testing.T, recordDir utilhub.FileNode, autoTestConfig utilhub.AutoConfigType) {
	for bpWidth := 0; bpWidth < len(autoTestConfig.Parameters.BpWidth); bpWidth++ {
		_runRandomizedBoundary(t, bpWidth, recordDir, autoTestConfig)
	}
}

// _runRandomizedBoundary 🧫 runs the actual test cases for RandomizedBoundary test.
func _runRandomizedBoundary(t *testing.T, bpWidth int, recordDir utilhub.FileNode, autoTestConfig utilhub.AutoConfigType) {
	dataChan, errChan, finishChan := recordDir.ReadBytesInChunksWithProgress("RandomizedBoundary.do_not_open", 8, binary.LittleEndian)

	root := NewBpTree(autoTestConfig.Parameters.BpWidth[bpWidth])

	testRandomizedBoundaryName := fmt.Sprintf("Randomized Boundary - run; Width: %3d", autoTestConfig.Parameters.BpWidth[bpWidth])

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		testRandomizedBoundaryName,
		// "RandomizedBoundary: Execution   ",                             // Progress bar title.
		uint32(autoTestConfig.Parameters.RandomTotalCount), // Total number of operations.
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
		case data := <-dataChan:
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
		case <-finishChan:
			break Loop
		}
	}

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Print a final report.
	err := progressBar.Report(len(testRandomizedBoundaryName + "; Width: XX"))
	assert.NoError(t, err)

	// Print the B Plus tree structure.
	root.root.Print()
}
