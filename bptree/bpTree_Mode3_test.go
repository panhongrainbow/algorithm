package bpTree

import (
	"fmt"
	"testing"

	bptestModel3 "github.com/panhongrainbow/algorithm/testdata/model3"

	bptestModel2 "github.com/panhongrainbow/algorithm/testdata/model2"
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

	bptest3 := &bptestModel3.BpTestModel3{}

	testDataSet, _ := bptest3.GenerateRandomSet()

	fmt.Println("长度", len(testDataSet))

	bptest2 := &bptestModel2.BpTestModel2{}

	err := bptest2.CheckRandomSet(testDataSet)

	fmt.Println("错误", err)

}

// verifyMode3 🧫 checks the test data set for Mode 3.
func verifyMode3(t *testing.T) {
	//
}

// runMode3 🧫 runs the actual test cases for Mode 3.
func runMode3(t *testing.T) {
	//
}

// _runMode3 🧫 runs the actual test cases for Mode 3.
func _runMode3(t *testing.T, bpWidth int) {
	//
}
