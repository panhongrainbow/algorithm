package utilhub

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

// LinuxSpliceStreamWrite wraps the LinuxSpliceStreamWrite function to write a file stream to a file.
func (fn FileNode) LinuxSpliceStreamWrite(filename string, fileFlag int, filePerm os.FileMode) (dataChan chan [][]byte, finishChan chan struct{}, err error) {
	// Construct the absolute path of the file by joining the transfer directory and the filename.
	absPath := path.Join(fn.transfer, filename)

	// Check if the directory exists.
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("file does not exist: %s", absPath)
	}

	// Check if the path is a file.
	if info.IsDir() {
		return nil, nil, fmt.Errorf("path is not a file: %s", absPath)
	}

	// Call the LinuxSpliceStreamWrite function with the absolute path and other parameters.
	return LinuxSpliceStreamWrite(absPath, fileFlag, filePerm)
}

// ReadBytesInChunks uses a goroutine to perform the file reading, allowing it to run concurrently with the main program flow.
func (fn FileNode) ReadBytesInChunks(filename string, chunkSize int) (<-chan []byte, <-chan error) {
	// Create channels to hold the chunked data and errors.
	dataChan := make(chan []byte)
	errChan := make(chan error, 1)

	// Construct the absolute path of the file by joining the transfer directory and the filename.
	absPath := path.Join(fn.transfer, filename)

	// Check if the directory exists.
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		// If the directory does not exist, send an error on the errChan and return.
		errChan <- fmt.Errorf("file does not exist: %s", absPath)
		return nil, errChan
	}

	// Check if the path is a file.
	if info.IsDir() {
		// If the path is not a file, send an error on the errChan and return.
		errChan <- fmt.Errorf("path is not a file: %s", absPath)
		return nil, errChan
	}

	// Start a goroutine to perform the file reading.
	go func() {
		// Defer closing the channels when the goroutine exits.
		defer close(dataChan)
		defer close(errChan)

		// Open the file for reading.
		file, err := os.Open(absPath)
		if err != nil {
			// If an error occurs while opening the file, send it on the errChan and return.
			errChan <- err
			return
		}
		// Defer closing the file when the goroutine exits.
		defer func() { _ = file.Close() }()

		// Create a buffered reader to read the file in chunks.
		reader := bufio.NewReader(file)
		// Create a buffer to hold the chunked data.
		buffer := make([]byte, chunkSize)

		// Read the file in chunks.
		for {
			// Read a chunk of data from the file.
			n, err := reader.Read(buffer)
			fmt.Println("n", n)
			// If data was read, send it on the dataChan.
			if n > 0 {
				dataChan <- buffer[:n]
			}
			// Check for errors.
			if err != nil && err.Error() == "EOF" {
				// If the end of the file was reached, send the error on the errChan and return.
				errChan <- err
				return
			}
			if err != nil && err.Error() != "EOF" {
				// If an error occurred while reading, send it on the errChan and return.
				errChan <- err
				return
			}
		}
	}()

	// Return the dataChan and errChan.
	return dataChan, errChan
}
