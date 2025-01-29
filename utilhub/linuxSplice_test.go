package utilhub

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"sync"
	"testing"
)

// Test_LinuxSpliceBulkWrite verifies the behavior of the LinuxSpliceBulkWrite function.
func Test_LinuxSpliceBulkWrite(t *testing.T) {
	// Skip this test if the operating system is not Linux, as the function is Linux-specific.
	if runtime.GOOS != "linux" {
		t.Skip("⏸️ Skipping test on non-Linux OS: " + t.Name())
	}

	// Define test cases for LinuxSpliceBulkWrite.
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
func Test_LinuxSpliceStreamWrite_Race(t *testing.T) {
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
			name:     "Feed Stream Data Continuously to LinuxSpliceStreamWrite by using goroutines",
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
			const iterations = 10
			// const iterations = 100000

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
