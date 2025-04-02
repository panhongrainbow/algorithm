package bpTree

import (
	"encoding/binary"
	"fmt"
	bptestModel1 "github.com/panhongrainbow/algorithm/testplan/bptestplan/model1"
	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func Test_Check_BpTree_Accuracy_mode1_preparation(t *testing.T) {

	// #################################################################################################
	// 🛠 Create some new instances for generating test data, showing progress, and writing test records.
	// #################################################################################################

	// -----> for Generating test data.

	// Create a new instance of BpTestModel1 with the specified random total count.
	bptest1 := &bptestModel1.BpTestModel1{RandomTotalCount: uint64(randomTotalCount)}

	// -----> for writing test records.

	// Create a new empty file named "mode0.do_not_open" under the record date path.
	err := recordDateNode.Touch("mode0.do_not_open")
	require.NoError(t, err, "record file could not be created; please check the path.")

	// #################################################################################################
	// 🛠 Generate and validate test data.
	// This is the most simple test case for testing the consistency and integrity of the B Plus Tree.
	// The half of the data is positive numbers and the other half is negative numbers.
	// #################################################################################################

	// Generate a random data set using the GenerateRandomSet method of BpTestModel1.
	// This method generates a slice of random data for testing purposes.
	// The half of the data is positive numbers and the other half is negative numbers.
	testDataSet, err := bptest1.GenerateRandomSet(1, 10)
	require.NoError(t, err, "test data set could not be generated; please check the parameters.")

	// Validate the generated data set using the CheckRandomSet method of BpTestModel1.
	// This method checks the validity of the data set by comparing the positive and negative numbers.
	err = bptest1.CheckRandomSet(testDataSet)
	require.NoError(t, err, "test data set could not be validated; please check the data set.")

	// #################################################################################################
	// 🛠 Initialize linux splice stream writer and set up some parameters for data writing.
	// #################################################################################################

	// -----> for Initializing linux splice.

	// Initialize a Linux splice stream writer to write data to a file.
	// The file is created with write-only permissions and truncated if it already exists.
	// spliceDataChan, spliceFinishChan, err := recordDateNode.LinuxSpliceStreamWrite("mode0.do_not_open", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// require.NoError(t, err, "Linux splice stream writer could not be initialized; please check the file path.")

	// -----> for setting up parameters.

	// Fixed Parameters:
	const (
		// Define constants for block length and width.
		// These constants determine the size of the blocks used for data writing.
		spliceBlockLength = 300 // The length of each block.
		spliceBlockWidth  = 100 // The width of each block.
	)

	// Variable Parameters:
	var (
	// spliceWritingPoint    = 0     // Initialize the start point for block writing.
	// spliceWritingFinished = false // Initialize a flag to track whether the writing process is finished.
	)

	// #################################################################################################
	// 🛠 Initialize linux splice stream writer and set up some parameters for data writing.
	// #################################################################################################

	err = recordDateNode.LinuxSpliceProgressStreamWrite(
		"test", utilhub.BrightCyan,
		"mode0.do_not_open", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
		testDataSet, binary.LittleEndian, spliceBlockLength, spliceBlockWidth)

	require.NoError(t, err)

	/*
		// Write data to the file in blocks until the entire data set is written.
		for !spliceWritingFinished {
			// Convert the data set to a block of bytes using the Int64SliceToBlockBytes method in utilhub.
			// This method converts a slice of int64 values to a block of bytes.
			var block [][]byte
			block, spliceWritingPoint, spliceWritingFinished, err = utilhub.Int64SliceToBlockBytes(testDataSet, binary.LittleEndian, spliceWritingPoint, spliceBlockLength, spliceBlockWidth)
			// Check if an error occurred during block writing.
			require.NoError(t, err)
			// Write the block to the file using the data channel.
			spliceDataChan <- block
		}

		// Close the data channel to signal the end of writing.
		close(spliceDataChan)

		// Wait for the finish channel to receive the finish signal.
		<-spliceFinishChan
	*/
}

func Test_Check_BpTree_Accuracy_mode1_execution(t *testing.T) {

	// #################################################################################################
	// 🛠 Decide the test method to execute.
	// Mode Identifier Number : 0
	// Mode Identifier Name   : Testing
	// Mode Description       : Make a bulk insert and bulk delete to test the consistency and integrity of the B Plus Tree.
	// #################################################################################################

	// Define a test mode for testing.
	testMode0Name := "Mode 0: Testing"

	// #################################################################################################
	// 🛠 The main test execution starts here.
	// #################################################################################################

	// Run the test mode.
	t.Run(testMode0Name, func(t *testing.T) {

		// #################################################################################################
		// 🛠
		// #################################################################################################

		var result []int64

		dataChan2, errChan2 := recordDateNode.ReadBytesInChunks("mode0.do_not_open", 800)
	Loop:
		for {
			select {
			case err := <-errChan2:
				if err == io.EOF {
					break Loop
				}
				if err != nil && err != io.EOF {
					assert.NoError(t, err)
				}
			case data := <-dataChan2:
				data2, _ := utilhub.BytesToInt64Slice(data, binary.LittleEndian)
				result = append(result, data2...)
				fmt.Println(len(result))
			}
		}

		fmt.Println(result[0])

		/*
			for i := 0; i < 10; i++ {
				data, err := utilhub.BytesToInt64Slice(<-dataChan2, binary.LittleEndian)
				require.NoError(t, err)
				result = append(result, data...)
			}
			fmt.Println(<-errChan2)
		*/

		// Check ...
		// content, _ := os.ReadFile("/home/tmp/" + time.Now().Format("2006-01-02") + "/mode0.do_not_open")
		// got, err := utilhub.BytesToInt64Slice(content[0:33], binary.LittleEndian)
		// fmt.Println(got, err)
		// fmt.Println(got[0:10], got[len(got)/2:len(got)/2+10])
		// fmt.Println(len(got))
		// err = bptest1.CheckRandomSet(got)
		// require.NoError(t, err)
	})
}
