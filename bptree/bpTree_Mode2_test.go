package bpTree

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	bptestModel2 "github.com/panhongrainbow/algorithm/testdata/model2"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy Mode 2 (Boundary Mode)
// Test cases are divided into three phases: preparation, validation, and execution.
// prepare_Mode1 : prepares test data for Mode 1.
// verify_Mode1 : validates the test data.
// run_Mode1 : runs the test cases.
// =====================================================================================================================

// prepareMode2 🧫 prepares test data for Mode 2.
func prepareMode2(t *testing.T) {

	// === Init test model and record file ===

	// Create model 2 with specified data count.
	bptest2 := &bptestModel2.BpTestModel2{}

	// Create an empty record file.
	err := recordDir.Touch("mode2.do_not_open")
	require.NoError(t, err, "failed to create record file")

	// === Generate test data ===

	// Generate a random set: half positive, half negative.
	testDataSet, err := bptest2.GenerateRandomSet()
	require.NoError(t, err, "failed to generate test data")

	// === Set write parameters ===

	const (
		spliceBlockLength = 300
		spliceBlockWidth  = 100
	)

	// === Write data with Linux splice stream ===

	fmt.Println("數量", len(testDataSet))

	err = recordDir.LinuxSpliceProgressStreamWrite(
		testDataSet,
		"mode2.do_not_open",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"Mode 2: Boundary - Backup",
		utilhub.BrightCyan,
		70,
	)
	require.NoError(t, err)

	// Data check is done in the next test case.
}
