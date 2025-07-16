package bpTree

import (
	"encoding/binary"
	"fmt"
	bptestModel1 "github.com/panhongrainbow/algorithm/testplan/bptestplan/model1"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// =====================================================================================================================
//                  ⚗️ BpTree Accuracy Mode 1 (Test Mode)
// These test cases are divided into three phases: preparation, Validating and execution.
// [Test_Check_BpTree_Accuracy_mode1_generation] is responsible for generating the test data for [Mode 1]. (Preparation 产生)
// [Test_Check_BpTree_Accuracy_mode1_check_test_data] verifies the integrity and correctness of the test data for [Mode 1]. (Validation 检查)
// [Test_Check_BpTree_Accuracy_mode1_execution] runs the actual test cases for [Mode 1]. (Execution 执行)
// =====================================================================================================================

// Test_Check_BpTree_Accuracy_mode1_generation 🧫 prepares the test data for [Mode 1].
func Test_Check_BpTree_Accuracy_mode1_generation(t *testing.T) {

	// #################################################################################################
	// Create some new instances for generating test data, showing progress, and writing test records. (初始化)
	// #################################################################################################

	// -----> for Generating test data.

	// Create a new instance of BpTestModel1 with the specified random total count.
	bptest1 := &bptestModel1.BpTestModel1{RandomTotalCount: uint64(bptreeUnitTestcfg.Parameters.RandomTotalCount)}

	// -----> for writing test records.

	// Create a new empty file named "mode0.do_not_open" under the record date path.
	err := recordDateNode.Touch("mode0.do_not_open")
	require.NoError(t, err, "record file could not be created; please check the path.")

	// #################################################################################################
	// Generate and validate test data. (产生测试数据)
	// This is the simplest test case for testing the consistency and integrity of the [B Plus Tree].
	// The half of the data is [positive numbers], and the other half is [negative numbers]. (一半为正，一半为负)
	// #################################################################################################

	// Generate a random data set using the [GenerateRandomSet] method of BpTestModel1.
	// This method generates a slice of random data for testing purposes.
	// The half of the data is positive numbers, and the other half is negative numbers.
	// 公式为: 最大值 max = 测试总数 total / 冲撞比例 collision_rate * 100 + 最小值 min
	// Then Check [MaxInt64] to avoid overflow.
	testDataSet, err := bptest1.GenerateRandomSet(1, 10)
	require.NoError(t, err, "test data set could not be generated; please check the parameters.")

	// #################################################################################################
	// Set up some parameters for [Linux Splicing] and Data Writing. (设定写入参数)
	// #################################################################################################

	// -----> for setting up parameters.

	// Fixed Parameters:
	// These parameters are related to writing speed. (在这里调速)
	const (
		// Define constants for block length and width.
		// These constants determine the size of the blocks used for data writing.
		spliceBlockLength = 300 // The length of each block.
		spliceBlockWidth  = 100 // The width of each block.
	)

	// #################################################################################################
	// Initialize [Linux Splice Stream Writer] and show progress. (开始 SPLICE 写入)
	// #################################################################################################

	err = recordDateNode.LinuxSpliceProgressStreamWrite(
		testDataSet,         // 原始数据，知道数据数量
		"mode0.do_not_open", // 要写入的文件名
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		binary.LittleEndian, spliceBlockLength, spliceBlockWidth,
		"Mode 1: Bulk Insert/Delete - Backup", // 进度条的标题
		utilhub.BrightCyan,                    // 进度条的颜色
		70,                                    // 进度条的显示长度
	)

	require.NoError(t, err)

	// Finish Checking the test data in the next test case. (下一个测试再檢查)
}

// Test_Check_BpTree_Accuracy_mode1_check_test_data 🧫 checks the test data set for [Mode 1].
func Test_Check_BpTree_Accuracy_mode1_check_test_data(t *testing.T) {

	testDataSet, err := recordDateNode.ReadAllBytesWithProgress(
		uint32(bptreeUnitTestcfg.Parameters.RandomTotalCount),
		"mode0.do_not_open", 800,
		binary.LittleEndian,
		"Mode 1: Bulk Insert/Delete - check Test Data", // 进度条的标题
		utilhub.BrightCyan,                             // 进度条的颜色
		70,                                             // 要读取的文件名
	)

	// -----> for Generating test data.

	// Create a new instance of BpTestModel1 with the specified random total count.
	bptest1 := &bptestModel1.BpTestModel1{RandomTotalCount: uint64(bptreeUnitTestcfg.Parameters.RandomTotalCount)}

	fmt.Println(len(testDataSet))

	// Validate the generated data set using the [CheckRandomSet] method of BpTestModel1.
	// This method checks the validity of the data set by comparing the positive and negative numbers.
	err = bptest1.CheckRandomSet(testDataSet)
	require.NoError(t, err, "test data set could not be validated; please check the data set.")
}

// Test_Check_BpTree_Accuracy_mode1_execution 🧫 runs the actual test cases for [Mode 1].
func Test_Check_BpTree_Accuracy_mode1_execution(t *testing.T) {
	dtatChan, errChan, finsishChan := recordDateNode.ReadBytesInChunksWithProgress("mode0.do_not_open", 8, binary.LittleEndian)

	root := NewBpTree(5)

	// ▓▒░ Creating a progress bar with optional configurations.
	progressBar, _ := utilhub.NewProgressBar(
		"Mode 1: Execution   ",                                // Progress bar title.
		uint32(bptreeUnitTestcfg.Parameters.RandomTotalCount), // Total number of operations.
		70,                                                    // Progress bar width.
		utilhub.WithTracking(5),                               // Update interval.
		utilhub.WithTimeZone("Asia/Taipei"),                   // Time zone.
		utilhub.WithTimeControl(500),                          // Update interval in milliseconds.
		utilhub.WithDisplay(utilhub.BrightGreen),              // Display style.
	)

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

Loop:
	for {
		select {
		case data := <-dtatChan:
			for i := 0; i < len(data); i++ {
				if data[i] >= 0 {
					root.InsertValue(BpItem{Key: data[i]})
					progressBar.UpdateBar()
				}
				if data[i] < 0 {
					deleted, _, _, err := root.RemoveValue(BpItem{Key: -1 * data[i]})
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

	root.root.Print()
}
