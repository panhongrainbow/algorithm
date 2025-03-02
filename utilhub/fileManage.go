package utilhub

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	dirPermission  = 0755
	filePermission = 0644
)

// FileManager represents a file manager that can create directories and track errors.
type FileManager struct {
	// transfer stores the current directory path being transferred.
	transfer string
	// err stores any errors that occur during file operations.
	err error
}

// MkDir creates a new directory at the specified path.
func (fm FileManager) MkDir(path string) FileManager {
	// Check if a previous error has occurred and return it if so.
	if fm.err != nil {
		fm.transfer = ""
		return fm
	}

	// Check if the directory already exists.
	if _, err := os.Stat(path); err == nil {
		// Directory already exists, return immediately without error.
		return fm
	}

	// Attempt to create the directory with the specified permissions.
	if err := os.MkdirAll(path, dirPermission); err != nil {
		// Return an error if directory creation fails.
		fm.err = fmt.Errorf("failed to create directory %s: %v", path, err)
		return fm
	}

	// Update the transfer path to include the newly created directory.
	fm.transfer = filepath.Join(fm.transfer, path)

	// Return nil to indicate successful directory creation.
	return fm
}

// Goto navigates to a specific directory path.
func (fm FileManager) Goto(path string) FileManager {
	// Check if a previous error has occurred and return it if so.
	if fm.err != nil {
		fm.transfer = ""
		return fm
	}

	// Check if the directory exists.
	// os.Stat returns a FileInfo describing the file, or an error if the file does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the directory does not exist, return an error with a descriptive message.
		fm.err = fmt.Errorf("directory does not exist: %s", path)
		return fm
	}

	// Update the transfer path to include the specified directory.
	// filepath.Join is used to join the current transfer path with the new path.
	fm.transfer = filepath.Join(fm.transfer, path)

	// Return nil to indicate successful navigation to the directory.
	return fm
}

// Jump navigates to a specific directory path by joining the provided paths to the current transfer path.
func (fm FileManager) Jump(paths ...string) FileManager {
	// Check if a previous error has occurred and return it if so.
	if fm.err != nil {
		fm.transfer = ""
		return fm
	}

	// Join the provided paths to the current transfer path using filepath.Join.
	fm.transfer = filepath.Join(fm.transfer, filepath.Join(paths...))

	// Check if the resulting directory exists using os.Stat.
	// If the directory does not exist, os.Stat returns an error.
	if _, err := os.Stat(fm.transfer); os.IsNotExist(err) {
		// If the directory does not exist, return an error with a descriptive message.
		fm.err = fmt.Errorf("directory does not exist: %s", fm.transfer)
		return fm
	}

	// If the directory exists, return the updated FileManager instance.
	return fm
}

// Touch creates a new empty file or truncates an existing file to zero length.
// If the file does not exist, it is created. If the file exists, its contents are cleared.
func (fm FileManager) Touch(filename string) error {
	// Check if a previous error has occurred and return it immediately.
	if fm.err != nil {
		fm.transfer = ""
		return fm.err
	}

	// Initialize the file path if it's not already set.
	if fm.transfer == "" {
		fm.transfer = filepath.Join("./", filename)
	} else {
		fm.transfer = filepath.Join(fm.transfer, filename)
	}

	// filename cannot be empty.
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check if the file exists.
	if _, err := os.Stat(fm.transfer); os.IsNotExist(err) {
		// File does not exist, create a new empty file.
		_, err := os.Create(fm.transfer)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", fm.transfer, err)
		}
	} else {
		// File exists, truncate its contents to zero length.
		err := os.WriteFile(fm.transfer, []byte(""), filePermission)
		if err != nil {
			return fmt.Errorf("failed to empty file %s: %v", fm.transfer, err)
		}
	}

	return nil
}

// FileTag represents a set of flags for selecting a date and time format.
type FileTag struct {
	// yearMonth indicates whether to include the year and month in the date format.
	yearMonth bool
	// yearMonthDay indicates whether to include the year, month, and day in the date format.
	yearMonthDay bool
	// yearMonthDayHour indicates whether to include the year, month, day, and hour in the date format.
	yearMonthDayHour bool
	// yearMonthDayHourMinute indicates whether to include the year, month, day, hour, and minute in the date format.
	yearMonthDayHourMinute bool
	// yearMonthDayHourMinuteSecond indicates whether to include the year, month, day, hour, minute, and second in the date format.
	yearMonthDayHourMinuteSecond bool
}

// DateTimeTag returns a formatted date and time string based on the provided FileTag.
// It returns an error if no valid time format is selected.
func DateTimeTag(ft FileTag) (string, error) {
	// Get the current time.
	now := time.Now()

	// Initialize the format string.
	var format string

	// Determine the date and time format based on the FileTag flags.
	switch {
	case ft.yearMonth:
		// Format: YYYY-MM
		format = "2006-01"
	case ft.yearMonthDay:
		// Format: YYYY-MM-DD
		format = "2006-01-02"
	case ft.yearMonthDayHour:
		// Format: YYYY-MM-DD HH
		format = "2006-01-02 15"
	case ft.yearMonthDayHourMinute:
		// Format: YYYY-MM-DD HH:MM
		format = "2006-01-02 15:04"
	case ft.yearMonthDayHourMinuteSecond:
		// Format: YYYY-MM-DD HH:MM:SS
		format = "2006-01-02 15:04:05"
	}

	// Check if a valid time format was selected.
	if format == "" {
		// Return an error if no valid time format was selected.
		return "", fmt.Errorf("no valid time format selected")
	}

	// Format the current time according to the selected format.
	timestamp := now.Format(format)

	// Return the formatted date and time string.
	return timestamp, nil
}

// List retrieves the directories and files inside a given directory.
// It returns three values: a slice of directories, a slice of files, and an error.
func (fm FileManager) List() (dir []string, file []string, err error) {
	// Check if a previous error has occurred and return it if so.
	if fm.err != nil {
		fm.transfer = ""
		return nil, nil, err
	}

	// Check if the given path exists and is a directory.
	var info os.FileInfo
	if info, err = os.Stat(fm.transfer); err != nil {
		return // Return empty slices and the error.
	}

	// Check if the path is a directory.
	if !info.IsDir() {
		err = fmt.Errorf("path is not a directory: %s", fm.transfer)
		return
	}

	// Open the directory.
	var dirHandle *os.File
	if dirHandle, err = os.Open(fm.transfer); err != nil {
		return // Return empty slices and the error.
	}
	defer dirHandle.Close() // Ensure the directory is closed.

	// Read directory contents.
	var entries []os.FileInfo
	if entries, err = dirHandle.Readdir(-1); err != nil {
		return // Return empty slices and the error.
	}

	// Iterate through directory entries.
	for _, entry := range entries {
		if entry.IsDir() {
			dir = append(dir, entry.Name()) // Append to directory slice.
		} else {
			file = append(file, entry.Name()) // Append to file slice.
		}
	}

	return // Implicit return of named values.
}

// validateAbsolutePath ensures that the given path is absolute.
func (fm FileManager) validateAbsolutePath(paths ...string) (string, error) {
	var absPath string
	for _, segment := range paths {
		absPath = filepath.Join(absPath, segment)
	}

	// Ensure the path is absolute.
	if !filepath.IsAbs(absPath) {
		return "", errors.New("parameter path is not absolute")
	}

	return absPath, nil
}

// RemoveFile removes the specified file if it exists and is an absolute path.
func (fm FileManager) RemoveFile(paths ...string) error {
	// Check the given path is absolute.
	absPath, err := fm.validateAbsolutePath(paths...)
	if err != nil {
		return err
	}

	// Get information about the file
	info, err := os.Stat(absPath)
	if err != nil {
		// If the file does not exist, return an error
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", absPath)
		}
		// If there is an error getting the file information, return the error
		return err
	}

	// Check if the path is a directory
	if info.IsDir() {
		// Return an error if the path is a directory, not a file
		return fmt.Errorf("path is a directory, not a file: %s", absPath)
	}

	// Attempt to remove the file
	if err := os.Remove(absPath); err != nil {
		// Return an error if the removal operation fails
		return fmt.Errorf("failed to remove file %s: %v", absPath, err)
	}

	// Return nil if the removal operation is successful
	return nil
}

// RemoveDir removes the specified directory if it exists and is an absolute path.
func (fm FileManager) RemoveDir(paths ...string) error {
	// Check the given path is absolute.
	absPath, err := fm.validateAbsolutePath(paths...)
	if err != nil {
		return err
	}

	// Check if the directory exists.
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}

	// Ensure the path is a directory.
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Attempt to remove the directory.
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("failed to remove directory %s: %v", absPath, err)
	}

	return nil
}
