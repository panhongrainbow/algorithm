package utilhub

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLinuxSpliceCopy verifies that a binary file with UUIDs can be copied correctly
// using LinuxSpliceBulkRead and LinuxSpliceBulkWrite.
func TestLinuxSpliceCopy(t *testing.T) {
	// Path to the source file.
	srcFile := "/tmp/test_source.bin"

	// Path to the destination file.
	dstFile := "/tmp/test_copy.bin"

	// Attempt to remove the source file and the destination file.
	defer func() {
		_ = os.Remove(srcFile)
		_ = os.Remove(dstFile)
	}()

	// Generate a binary file with 10 random UUIDs.
	var err error
	err = GenerateBinaryFile(srcFile, 10)
	require.Nil(t, err, "GenerateBinaryFile failed")

	// Read the source file into chunks.
	var chunks [][]byte
	chunks, err = LinuxSpliceBulkRead(srcFile, 32) // 32 bytes per chunk
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
