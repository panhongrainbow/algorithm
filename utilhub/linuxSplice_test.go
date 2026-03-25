package utilhub

import (
	"bytes"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLinuxSpliceCopy verifies that a binary file with UUIDs can be copied correctly
// using LinuxSpliceBulkRead and LinuxSpliceBulkWrite.
func TestLinuxSpliceCopy(t *testing.T) {
	srcFile := "/tmp/test_source.bin"
	dstFile := "/tmp/test_copy.bin"
	// defer func() { _ = os.Remove(srcFile) }()
	// defer func() { _ = os.Remove(dstFile) }()

	// Generate a binary file with 10 random UUIDs.
	var err error
	err = GenerateBinaryFile(srcFile, 10)
	require.Nil(t, err, "GenerateBinaryFile failed")

	// Read the source file into chunks.
	chunks, err := LinuxSpliceBulkRead(srcFile, 32) // 32 bytes per chunk
	require.Nil(t, err, "LinuxSpliceBulkRead failed")

	// Write the chunks to the destination file.
	if err = LinuxSpliceBulkWrite(dstFile, chunks, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		t.Fatalf("Failed to write copy file: %v", err)
	}

	// Read the copied file into chunks.
	newChunks, err := LinuxSpliceBulkRead(dstFile, 32)
	require.Nil(t, err, "LinuxSpliceBulkRead failed")

	// Compare each UUID in order.
	originalData := bytes.Join(chunks, nil)
	copiedData := bytes.Join(newChunks, nil)

	// Verify that the copied file matches the original and both file sizes are multiples of UUID size (16 bytes).
	require.Equal(t, originalData, copiedData, "Copied file does not match original")
	require.Equal(t, 0, len(originalData)%16, "Original file size is not multiple of UUID size")
	require.Equal(t, 0, len(copiedData)%16, "Copied file size is not multiple of UUID size")

	// Compare each 16-byte UUID in the original and copied data to ensure every UUID matches.
	for i := 0; i < len(originalData); i += 16 {
		require.True(t, bytes.Equal(originalData[i:i+16], copiedData[i:i+16]))
	}
}

// Test_LinuxSpliceBulkWrite verifies that bulk writes are correctly performed using Linux splice.
func Test_LinuxSpliceBulkWrite(t *testing.T) {
	// This test is skipped on non-Linux platforms, as LinuxSpliceBulkWrite is Linux-specific.
	if runtime.GOOS != "linux" {
		t.Skip("⏸️ Skipping test on non-Linux OS: " + t.Name())
	}

	// Defines the set of test cases for LinuxSpliceBulkWrite.
	tests := []struct {
		name        string      // Descriptive name of the test case
		filename    string      // File path to write to
		data        [][]byte    // Data chunks to write
		fileFlag    int         // Flags used when opening the file (e.g., create, write-only, truncate)
		filePerm    os.FileMode // Permissions for the created file
		wantErr     bool        // Indicates if an error is expected
		wantContent string      // Expected content to be written to the file
	}{
		{
			// Test case: Successful bulk write to the file.
			name:        "Simple Test for LinuxSpliceBulkWrite",
			filename:    "/tmp/test_file.txt",
			data:        [][]byte{[]byte("Hello"), []byte(" "), []byte("World"), []byte("!")},
			fileFlag:    os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
			filePerm:    0644,
			wantErr:     false,
			wantContent: "Hello World!",
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test with the specified parameters.
			err := LinuxSpliceBulkWrite(tt.filename, tt.data, tt.fileFlag, tt.filePerm)

			// Verify if the error behavior matches the expectation.
			if tt.wantErr {
				// If an error is expected, assert that one was returned.
				assert.Error(t, err)
			} else {
				// If no error is expected, assert the function completed successfully.
				assert.NoError(t, err)

				// Read the content from the file to validate the result.
				content, err := os.ReadFile(tt.filename)
				assert.NoError(t, err)

				// Ensure the file content matches the expected result.
				assert.Equal(t, tt.wantContent, string(content))
			}

			// Clean up: Remove the test file after each test case.
			_ = os.Remove(tt.filename)
		})
	}
}

// Test_LinuxSpliceStreamWrite validates the behavior of the LinuxSpliceStreamWrite function.
func Test_LinuxSpliceStreamWrite(t *testing.T) {
	// Skip the test if the operating system is not Linux, as the function is Linux-specific.
	if runtime.GOOS != "linux" {
		t.Skip("⏸️ Skipping test on non-Linux OS: " + t.Name())
	}

	// Define test cases with various parameters and expected outcomes.
	tests := []struct {
		name        string      // Descriptive name of the test case.
		filename    string      // Path of the file to write to.
		data        [][]byte    // Data chunks to be written to the file.
		fileFlag    int         // File flags for opening the file (e.g., create, write-only, truncate).
		filePerm    os.FileMode // File permissions for the created file.
		wantErr     bool        // Whether an error is expected from the function.
		wantContent string      // Expected content to be written to the file.
	}{
		{
			// Test case: Successfully writing data to a file using stream.
			name:        "Simple Test for LinuxSpliceStreamWrite",
			filename:    "/tmp/test_file.txt",
			data:        [][]byte{[]byte("Hello"), []byte(" "), []byte("World"), []byte("!")},
			fileFlag:    os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
			filePerm:    0644,
			wantErr:     false,
			wantContent: "Hello World!",
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Attempt to initialize the LinuxSpliceStreamWrite function.
			// The function returns:
			// - A data channel to send data chunks for writing.
			// - A finish channel to signal when the write operation is complete.
			// - An error, if initialization fails.
			dataChan, finishChan, err := LinuxSpliceStreamWrite(tt.filename, tt.fileFlag, tt.filePerm)

			// Verify if the error behavior matches the expectation.
			if tt.wantErr {
				// If an error is expected, assert that one was returned.
				assert.Error(t, err)
			} else {
				// If no error is expected, ensure the initialization was successful.
				assert.NoError(t, err)

				// Send data chunks to the data channel for writing.
				dataChan <- tt.data

				// Close the data channel to indicate no more data will be sent.
				close(dataChan)

				// Wait for the write operation to complete.
				// This mechanism ensures all data has been written before proceeding.
				<-finishChan

				// Read the file content to validate the written data.
				content, err := os.ReadFile(tt.filename)
				require.NoError(t, err)

				// Assert that the file content matches the expected result.
				assert.Equal(t, tt.wantContent, string(content))

				// Clean up: Remove the test file after each test case.
				_ = os.Remove(tt.filename)
			}
		})
	}
}

// Test_LinuxSpliceStreamWrite_FeedStreamData validates the LinuxSpliceStreamWrite function by continuously writing data in batches
// and verifying the file content matches the expected pattern.
func Test_LinuxSpliceStreamWrite_FeedStreamData(t *testing.T) {
	// Skip the test if the operating system is not Linux, as the function is Linux-specific.
	if runtime.GOOS != "linux" {
		t.Skip("⏸️ Skipping test on non-Linux OS: " + t.Name())
	}

	// Define the test cases for the LinuxSpliceStreamWrite function.
	tests := []struct {
		name     string      // The name of the test case.
		filename string      // The path of the file to write data to.
		fileFlag int         // The file flags used for opening the file.
		filePerm os.FileMode // The file permissions applied when creating the file.
		wantErr  bool        // Indicates whether an error is expected during execution.
	}{
		{
			name:     "Feed Stream Data Continuously to LinuxSpliceStreamWrite",
			filename: "/tmp/test_file.txt",
			fileFlag: os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
			filePerm: 0644,
			wantErr:  false,
		},
	}

	// Iterate through the defined test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the splice stream writer with the specified file configuration.
			dataChan, finishChan, err := LinuxSpliceStreamWrite(tt.filename, tt.fileFlag, tt.filePerm)
			assert.NoError(t, err)

			// Define the number of iterations for writing and the size of each data batch.
			const iterations = 1000000
			const batchSize = 32 // The number of bytes in each batch.

			// Expected file size calculation.
			// Total file size = iterations * batchSize = 1,000,000 * 32 bytes = 32,000,000 bytes (30.52 MB).

			// Dynamically write data in batches using a predefined pattern.
			for i := 0; i < iterations; i++ {
				batch := make([]byte, batchSize)
				for j := 0; j < batchSize; j++ {
					// Generate data using a sequence of numbers modulo 256.
					batch[j] = byte((i*batchSize + j) % 256)
				}
				// Send the generated batch to the data channel.
				dataChan <- [][]byte{batch}
			}

			// Close the data channel to signal that no more data will be sent.
			close(dataChan)

			// Wait for the write operation to complete.
			<-finishChan

			// Read the content of the written file.
			content, err := os.ReadFile(tt.filename)
			require.NoError(t, err)

			// Generate the expected content based on the same dynamic pattern.
			expectedContent := make([]byte, iterations*batchSize)
			for i := range expectedContent {
				expectedContent[i] = byte(i % 256)
			}

			// Verify that the actual file content matches the expected pattern.
			assert.Equal(t, expectedContent, content)

			// Clean up: Remove the test file after each test case.
			_ = os.Remove(tt.filename)
		})
	}
}

// Test_LinuxSpliceStreamWrite_Race verifies that LinuxSpliceStreamWrite correctly writes ASCII data (0–255) to a file using concurrent goroutines.
// It ensures the file content matches the expected pattern and cleans up after the test.

// When testing for race conditions in Go, there are two main ways to verify whether a race condition is happening :

// 1. Using Go's Race Detector
// Go provides a built-in race detector that can identify when multiple goroutines access shared memory concurrently in an unsafe manner
// $ Test_LinuxSpliceStreamWrite_Race_Goroutines_Number=20 go test -race -v -run Test_LinuxSpliceStreamWrite_Race

// 2. Verifying Data Integrity
// Even if the Go race detector doesn't show any issues,
// I can still check for data integrity manually by verifying that the data written to the file or memory is consistent and correct after concurrent writes.
// $ Test_LinuxSpliceStreamWrite_Race_Goroutines_Number=100000 go test -v -run Test_LinuxSpliceStreamWrite_Race
func Test_LinuxSpliceStreamWrite_Race(t *testing.T) {
	// Skip the test if the operating system is not Linux, as the function is Linux-specific.
	if runtime.GOOS != "linux" {
		t.Skip("⏸️ Skipping test on non-Linux OS: " + t.Name())
	}

	// Define the test cases for the LinuxSpliceStreamWrite function.
	tests := []struct {
		name     string      // The name of the test case.
		filename string      // The name of the file to which data will be written.
		fileFlag int         // Flags for opening the file, such as read/write permissions.
		filePerm os.FileMode // The file permissions for the newly created file.
		wantErr  bool        // Indicates whether an error is expected during the test.
	}{
		{
			// Test case: Continuously feed data using ASCII codes.
			name:     "Feed Stream Data Continuously to LinuxSpliceStreamWrite by using Goroutines",
			filename: "/tmp/test_file.txt",                   // Temporary file for testing.
			fileFlag: os.O_CREATE | os.O_WRONLY | os.O_TRUNC, // File will be created, write-only, and truncated if exists.
			filePerm: 0644,                                   // Standard file permissions.
			wantErr:  false,                                  // No error expected.
		},
	}

	// Iterate through the defined test cases.
	for _, tt := range tests {
		// Run the test case defined by the test name.
		t.Run(tt.name, func(t *testing.T) {

			// Call the function under test to initialize the channels for writing.
			// The function returns data and finish channels, along with an error if any occurs.
			dataChan, finishChan, err := LinuxSpliceStreamWrite(tt.filename, tt.fileFlag, tt.filePerm)
			assert.NoError(t, err) // Assert that no error occurred while setting up the write.

			// Define the number of iterations for the test and the size of each batch of data to write.
			iterations := 10

			// The number of GoRoutines will only be increased if strict conditions are met;
			// otherwise, the original value will be maintained to ensure the system operates safely.

			// Get the value of the environment variable "Test_LinuxSpliceStreamWrite_Race_Goroutines_Number".
			// This variable determines the number of goroutines to be used for writing data.
			envVar := os.Getenv("Test_LinuxSpliceStreamWrite_Race_Goroutines_Number")

			// Convert the environment variable value to an integer.
			specificGoroutines, err := strconv.Atoi(envVar)
			if err == nil && envVar != "" && specificGoroutines > 10 {
				// If the conversion is successful and the value is greater than 10, set the number of goroutines to 100000.
				// Otherwise, keep the original value.
				iterations = specificGoroutines
			}

			// Calculate the expected total file size.
			// The file size is the number of iterations (100,000) multiplied by the batch size (256 bytes).
			// Total file size = 100,000 * 256 bytes = 25,600,000 bytes (approx. 24.41 MB).

			// Create a WaitGroup to ensure all goroutines complete their work before closing the test.
			var wg sync.WaitGroup
			wg.Add(iterations) // Add the number of iterations to the wait group.

			// Launch a goroutine for each iteration to write data.
			for i := 0; i < iterations; i++ {
				go func() {
					// Create a slice to hold the data to be written in each batch.
					eachData := make([][]byte, 0)

					// Write ASCII codes from 0 to 255 for each batch of data.
					// This loop generates a single byte for every ASCII code.
					for j := 0; j < 256; j++ {
						batch := make([]byte, 1)           // Create a batch with one byte.
						batch[0] = byte(j)                 // Assign the ASCII code to the byte.
						eachData = append(eachData, batch) // Append the byte to the data slice.
					}

					// Send the batch of data to the data channel for writing.
					dataChan <- eachData

					// Decrement the wait group counter as the goroutine completes.
					wg.Done()
				}()
			}

			// Wait for all goroutines to finish before proceeding.
			wg.Wait()

			// Close the data channel to indicate no more data will be sent.
			close(dataChan)

			// Wait for the write operation to be completed by the finish channel.
			<-finishChan

			// Read the file content after the write operation is complete.
			content, err := os.ReadFile(tt.filename)
			require.NoError(t, err) // Assert that no error occurred during file reading.

			// Generate the expected content for the file based on the ASCII pattern.
			expectedContent := make([]byte, 256*iterations)
			for i := range expectedContent {
				expectedContent[i] = byte(i % 256) // Use ASCII code in a cyclic pattern.
			}

			// Assert that the content read from the file matches the expected pattern.
			assert.Equal(t, expectedContent, content)

			// Clean up: Remove the test file after the test case has completed to maintain a clean environment.
			_ = os.Remove(tt.filename)
		})
	}
}
