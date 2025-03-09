package utilhub

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test_FileManager_MkDir tests the MkDir method of the FileManager struct.
func Test_FileManager_MkDir(t *testing.T) {
	// Create a new instance of the FileManager struct.
	fm := FileManager{}

	// Generate a unique directory name using uuid.New().String().
	dirName := "/tmp/" + uuid.New().String()

	// Test case: Create a new directory
	// Verify that no error is returned when creating a new directory.
	fm.MkDir(dirName)
	assert.NoError(t, fm.err)

	// Test case: Create a directory that already exists
	// Verify that no error is returned when creating a directory that already exists.
	fm.MkDir(dirName)
	assert.NoError(t, fm.err)

	// Test case: Create a directory with an error
	// Set the error field of the FileManager struct to "previous error".
	fm.err = fmt.Errorf("previous error")
	// Verify that the error returned when creating a directory is "previous error".
	fm.MkDir(dirName)
	assert.EqualError(t, fm.err, "previous error")
}

// Test_FileManager_Jump tests the Jump method of the FileManager struct.
func Test_FileManager_Jump(t *testing.T) {
	// Create an empty FileManager instance.
	nodeEmpty := FileManager{}

	// Navigate to the "/tmp" directory.
	nodeTmp := nodeEmpty.Goto("/tmp")

	// Generate a unique directory name for the first layer.
	layer1 := uuid.New().String()
	// Create a new directory at the first layer.
	nodeLayer1 := nodeTmp.MkDir(layer1)

	// Generate a unique directory name for the second layer.
	layer2 := uuid.New().String()
	// Create a new directory at the second layer.
	nodeLayer2 := nodeLayer1.MkDir(layer2)

	// Jump to the specific directory path using the Jump method.
	nodeJump := nodeEmpty.Jump("/tmp", layer1, layer2)

	// Verify that the transfer path of nodeJump is equal to the transfer path of nodeLayer2.
	assert.Equal(t, nodeLayer2.transfer, nodeJump.transfer)
}

// Test_FileManager_Goto tests the Goto method of the FileManager struct.
func Test_FileManager_Goto(t *testing.T) {
	// Create a new instance of the FileManager struct.
	fm := FileManager{}

	// Generate a unique directory name using uuid.New().String().
	dirName := "/tmp/" + uuid.New().String()

	// Test case: Create a new directory
	// Verify that no error is returned when creating a new directory.
	fm.MkDir(dirName)
	assert.NoError(t, fm.err)

	// Test case: Navigate to a directory that exists.
	fm.Goto(dirName)
	assert.NoError(t, fm.err)

	// Test case: Navigate to a directory that doesn't exist.
	err := fm.Goto("/tmp/" + uuid.New().String()).Error()
	assert.ErrorContains(t, err, "directory does not exist:")

	// Test case: Navigate to a directory with an error.
	fm.err = fmt.Errorf("previous error")
	fm.Goto(dirName)
	assert.EqualError(t, fm.err, "previous error")
}

// Test_FileManager_Touch tests the Touch method of the FileManager struct.
func Test_FileManager_Touch(t *testing.T) {
	// Generate a unique filename and directory name for testing.
	fileName := uuid.New().String() + ".txt"
	dirName := "/tmp/" + uuid.New().String()

	// Define test cases for the Touch method.
	tests := []struct {
		name        string
		fileManager FileManager
		dirname     string
		filename    string
		expectedErr string
	}{
		{
			// Test case: No filename provided.
			name: "no filename provided",
			fileManager: FileManager{
				transfer: "",
			},
			dirname:     "/tmp",
			filename:    "",
			expectedErr: "filename cannot be empty",
		},
		{
			// Test case: Create new file in existing directory.
			name: "create new file in existing directory",
			fileManager: FileManager{
				transfer: "",
			},
			dirname:     "/tmp",
			filename:    fileName,
			expectedErr: "",
		},
		{
			// Test case: Create new file in non-existing directory.
			name: "create new file in non-existing directory",
			fileManager: FileManager{
				transfer: "",
			},
			dirname:     dirName,
			filename:    fileName,
			expectedErr: "directory does not exist",
		},
		{
			// Test case: Return previous error.
			name: "return previous error",
			fileManager: FileManager{
				err: fmt.Errorf("previous error"),
			},
			dirname:     dirName,
			filename:    fileName,
			expectedErr: "previous error",
		},
	}

	// Run each test case.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the Touch method and store the error.
			err := test.fileManager.Goto(test.dirname).Touch(test.filename)

			// Check if an error is expected.
			if test.expectedErr != "" {
				// Assert that the error contains the expected error message.
				assert.ErrorContains(t, err, test.expectedErr)
			} else {
				// Assert that no error occurred.
				assert.NoError(t, err)
			}
		})
	}
}

// Test_DateTimeTag tests the DateTimeTag function for different FileTag configurations.
// TestDateTimeTag tests the DateTimeTag function for different FileTag configurations.
func Test_DateTimeTag(t *testing.T) {
	// Helper function to test the formatted date and time
	tests := []struct {
		name   string
		ft     FileTag
		format string
	}{
		{
			name:   "yearMonth",
			ft:     FileTag{yearMonth: true},
			format: "2006-01",
		},
		{
			name:   "yearMonthDay",
			ft:     FileTag{yearMonthDay: true},
			format: "2006-01-02",
		},
		{
			name:   "yearMonthDayHour",
			ft:     FileTag{yearMonthDayHour: true},
			format: "2006-01-02 15",
		},
		{
			name:   "yearMonthDayHourMinute",
			ft:     FileTag{yearMonthDayHourMinute: true},
			format: "2006-01-02 15:04",
		},
		{
			name:   "yearMonthDayHourMinuteSecond",
			ft:     FileTag{yearMonthDayHourMinuteSecond: true},
			format: "2006-01-02 15:04:05",
		},
		{
			name:   "noFormatSelected",
			ft:     FileTag{}, // No format selected
			format: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call DateTimeTag to get the formatted timestamp
			result, err := DateTimeTag(tt.ft)

			// If no format is selected, we expect an error
			if tt.format == "" {
				require.Error(t, err, "expected error for empty format, got nil")
				return
			}

			// If a format is selected, check if the result matches the expected format
			require.NoError(t, err, "unexpected error")

			// Now check if the result matches the expected format
			_, err = time.Parse(tt.format, result)
			assert.NoError(t, err, "expected format %s, but got %s, error: %v", tt.format, result, err)
		})
	}
}

// Test_FileManager_validateAbsolutePath tests the validateAbsolutePath method of the FileManager struct.
// This test ensures that the method correctly handles absolute and relative paths, as well as edge cases such as empty paths.
func Test_FileManager_validateAbsolutePath(t *testing.T) {
	// Define a slice of test cases, each representing a different scenario to test.
	tests := []struct {
		name       string   // The name of the test case.
		paths      []string // The input paths to test.
		wantErr    bool     // Whether an error is expected.
		expectPath string   // The expected absolute path.
	}{
		{
			name:       "absolute path 1/2",
			paths:      []string{"/", "path", "to", "file"},
			wantErr:    false,
			expectPath: "/path/to/file",
		},
		{
			name:       "absolute path 2/2",
			paths:      []string{"/////", "///path///", "///to///", "///file"},
			wantErr:    false,
			expectPath: "/path/to/file",
		},
		{
			name:       "relative path",
			paths:      []string{"path", "to", "file"},
			wantErr:    true,
			expectPath: "",
		},
		{
			name:       "empty path",
			paths:      []string{},
			wantErr:    true,
			expectPath: "",
		},
		{
			name:       "multiple absolute paths",
			paths:      []string{"/path/to/file", "/another/path"},
			wantErr:    false,
			expectPath: "/path/to/file/another/path",
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run the test case with a descriptive name.
		t.Run(tt.name, func(t *testing.T) {
			// Create a new FileManager instance for the test.
			fm := &FileManager{}

			// Call the validateAbsolutePath method with the input paths.
			absPath, err := fm.validateAbsolutePath(tt.paths...)

			// If an error is expected, check that one was returned.
			if tt.wantErr {
				require.Error(t, err)
			} else {
				// If no error is expected, check that none was returned.
				require.NoError(t, err)
				// Check that the returned absolute path is not nil.
				require.NotNil(t, absPath)
				// Check that the returned absolute path matches the expected path.
				require.Equal(t, tt.expectPath, absPath)
			}
		})
	}
}

// Test_FileManager_RemoveFile tests the RemoveFile method of the FileManager struct.
// This test case ensures that a file can be successfully removed after creation.
func Test_FileManager_RemoveFile(t *testing.T) {
	// Generate a unique filename and directory name for testing.
	dirName := "/tmp/" + uuid.New().String()
	fileName := uuid.New().String() + ".txt"

	// Create a new instance of the FileManager struct.
	fm := FileManager{}

	// Create a new directory and file for testing.
	err := fm.MkDir(dirName).Touch(fileName)
	// Assert that no error occurred during directory and file creation.
	assert.NoError(t, err, "expected no error, got %v", err)

	// Attempt to remove the file.
	err = fm.RemoveFile(dirName, fileName)
	// Assert that no error occurred during file removal.
	assert.NoError(t, err, "expected no error, got %v", err)

	return
}

// Test_FileManager_RemoveDir tests the RemoveDir method of the FileManager struct.
// This test case ensures that a directory can be successfully removed after creation.
func Test_FileManager_RemoveDir(t *testing.T) {
	// Generate a unique directory name and two unique file names for testing.
	dirName := "/tmp/" + uuid.New().String()
	fileName1 := uuid.New().String() + ".txt"
	fileName2 := uuid.New().String() + ".txt"

	// Print the generated directory and file names for debugging purposes.
	fmt.Println(dirName, fileName1, fileName2)

	// Create a new instance of the FileManager struct.
	fm := FileManager{}

	// Create a new directory and two files for testing.
	err := fm.MkDir(dirName).Touch(fileName1)
	// Assert that no error occurred during directory and file creation.
	assert.NoError(t, err, "expected no error, got %v", err)

	err = fm.MkDir(dirName).Touch(fileName2)
	// Assert that no error occurred during directory and file creation.
	assert.NoError(t, err, "expected no error, got %v", err)

	// Verify that the directory contains the two created files.
	_, files, err := fm.Goto(dirName).List()
	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, 2, len(files), "expected 2 files, got %v", len(files))

	// Attempt to remove the directory.
	err = fm.RemoveDir(dirName)
	assert.NoError(t, err, "expected no error, got %v", err)

	// Verify that the directory is empty after removal.
	_, files, err = fm.Goto(dirName).List()
	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, 0, len(files), "expected 0 files, got %v", len(files))

	return
}
