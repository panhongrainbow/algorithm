package bpTree

import (
	"encoding/binary"
	"os"
	"testing"

	bptestModel3 "github.com/panhongrainbow/algorithm/testdata/model3"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy Mode 3 (Endurance Test Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// prepare_Mode3 : prepares test data for Mode 3.
// verify_Mode3 : validates the test data.
// run_Mode3 : runs the test cases.
// =====================================================================================================================

// prepareMode3 🧫 prepares test data for Mode 3.
func prepareMode3(t *testing.T) {

	// === Init test model and record file ===

	// Create model 3 with specified data count.
	bptest3 := &bptestModel3.BpTestModel3{}

	// Create an empty record file.
	err := recordDir.Touch("mode3.do_not_open")
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
		"mode3.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"Mode 3: cyclicStress - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}

// verifyMode3 🧫 checks the test data set for Mode 3.
func verifyMode3(t *testing.T) {
	// Read test data with progress bar.
	testDataSet, err := recordDir.ReadAllBytesWithProgress(
		uint32(unitTestConfig.Parameters.RandomTotalCount),
		"mode3.do_not_open", 800,
		binary.LittleEndian,
		"Mode 3: Randomized Boundary Test - read test data",
		utilhub.BrightCyan,
		70,
	)

	// Init test model.
	bptest3 := &bptestModel3.BpTestModel3{}

	// Validate test data.
	err = bptest3.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to validate test data")
}

// runMode3 🧫 runs the actual test cases for Mode 3.
func runMode3(t *testing.T) {
	//
}
