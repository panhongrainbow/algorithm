package bpTree

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	bptestSingleNodeEndurance "github.com/panhongrainbow/go-algorithm/testdata/singleNodeEndurance"
	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy SingleNodeEndurance (Endurance Test Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// prepare_SingleNodeEndurance : prepares test data for SingleNodeEndurance.
// verify_SingleNodeEndurance : validates the test data.
// run_SingleNodeEndurance : runs the test cases.
// =====================================================================================================================

// prepareSingleNodeEndurance 🧫 prepares test data for SingleNodeEndurance.
func prepareMode3(t *testing.T) {

	// === Init test model and record file ===

	// Create SingleNodeEndurance with specified data count.
	bptest3 := &bptestSingleNodeEndurance.BpTestSingleNodeEndurance{}

	// Create an empty record file.
	err := recordDir.Touch("SingleNodeEndurance.do_not_open")
	require.NoError(t, err, "failed to create record file")

	// === Generate test data ===

	// Generate metal fatigue test data
	testDataSet, err := bptest3.GenerateRandomSet()
	require.NoError(t, err, "failed to generate test data")

	// === Set write parameters ===

	const (
		spliceBlockLength = 300
		spliceBlockWidth  = 100
	)

	// === Write data with Linux splice stream ===

	err = recordDir.LinuxSpliceProgressStreamWrite(
		testDataSet,
		"SingleNodeEndurance.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"SingleNodeEndurance - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}

// verifySingleNodeEndurance 🧫 checks the test data set for SingleNodeEndurance test.
func verifySingleNodeEndurance(t *testing.T) {
	// Read test data with progress bar.
	testDataSet, err := recordDir.ReadAllBytesWithProgress(
		uint32(autoTestConfig.Parameters.RandomTotalCount),
		"SingleNodeEndurance.do_not_open", 800,
		binary.LittleEndian,
		"SingleNodeEndurance - read test data",
		utilhub.BrightCyan,
		70,
	)

	// Init test model.
	bptest3 := &bptestSingleNodeEndurance.BpTestSingleNodeEndurance{}

	// Validate test data.
	err = bptest3.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to validate test data")
}

// runSingleNodeEndurance 🧫 runs the actual test cases for SingleNodeEndurance test.
func runSingleNodeEndurance(t *testing.T) {
	for bpWidth := 0; bpWidth < len(autoTestConfig.Parameters.BpWidth); bpWidth++ {
		_runSingleNodeEndurance(t, bpWidth)
	}
}

// _runSingleNodeEndurance 🧫 runs the actual test cases for SingleNodeEndurance test.
func _runSingleNodeEndurance(t *testing.T, bpWidth int) {
	dataChan, errChan, finishChan := recordDir.ReadBytesInChunksWithProgress("SingleNodeEndurance.do_not_open", 8, binary.LittleEndian)

	root := NewBpTree(autoTestConfig.Parameters.BpWidth[bpWidth])

	testSingleNodeEnduranceName := fmt.Sprintf("SingleNodeEndurance - run; Width: %3d", autoTestConfig.Parameters.BpWidth[bpWidth])

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		testSingleNodeEnduranceName,
		// "SingleNodeEndurance: Execution   ",                             // Progress bar title.
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
	err := progressBar.Report(len(testSingleNodeEnduranceName + "; Width: XX"))
	assert.NoError(t, err)

	// Print the B Plus tree structure.
	root.root.Print()
}
