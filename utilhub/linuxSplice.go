package utilhub

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// =====================================================================================================================
//                  🛠️ Linux Splice (Tool)
// Linux Splice is a tool for writing data to a file using the Splice system call.
// It is because there are a lot of data that need to be written to a file,
// so it is necessary to use the Splice system call to write data to the file.
// =====================================================================================================================

// LinuxSpliceBulkWrite ⛏️ writes multiple chunks of data to a file using the Splice system call.
// It creates a pipe to write the data to the pipe, and then uses the Splice system call to write the data to the file.
func LinuxSpliceBulkWrite(filename string, data [][]byte, fileFlag int, filePerm os.FileMode) error {
	// The file is opened with the specified flags and permissions.
	file, err := os.OpenFile(filename, fileFlag, filePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// A pipe is created to write the data to.
	pipe := make([]int, 2)
	if err := syscall.Pipe(pipe); err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	defer func() { _ = syscall.Close(pipe[0]) }()
	defer func() { _ = syscall.Close(pipe[1]) }()

	// Each chunk of data is written to the pipe.
	for _, chunk := range data {
		n, err := syscall.Write(pipe[1], chunk)
		if err != nil {
			return fmt.Errorf("failed to write to pipe: %w", err)
		}
		if n != len(chunk) {
			return fmt.Errorf("partial write to pipe, wrote %d bytes out of %d", n, len(chunk))
		}

		// The data is then written from the pipe to the file using the Splice system call.
		for n > 0 {
			written, err := syscall.Splice(pipe[0], nil, int(file.Fd()), nil, n, 0)
			if err != nil {
				return fmt.Errorf("failed to splice data: %w", err)
			}
			n -= int(written)
		}
	}

	return nil
}

// LinuxSpliceStreamWrite ⛏️ creates a pipe to write data to a file using the Splice system call.
// It returns a channel to send data to be written to the file.
func LinuxSpliceStreamWrite(filename string, fileFlag int, filePerm os.FileMode) (dataChan chan [][]byte, finishChan chan struct{}, err error) {
	// Open the file with the specified flags and permissions.
	file, err := os.OpenFile(filename, fileFlag, filePerm)
	if err != nil {
		// If the file cannot be opened, return an error.
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Create a pipe to write data to.
	pipe := make([]int, 2)
	if err := syscall.Pipe(pipe); err != nil {
		// If the pipe cannot be created, close the file and return an error.
		_ = file.Close()
		return nil, nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	// Create a channel to send data to be written to the file.
	dataChan = make(chan [][]byte, 100)

	finishChan = make(chan struct{})

	// Start a goroutine to write data to the file.
	go func() {
		// Defer closing the pipe and file.
		defer func() {
			// Close the pipe.
			_ = syscall.Close(pipe[0])
			_ = syscall.Close(pipe[1])

			// Close the file with retries.
			for i := 0; i < 5; i++ {
				if err := file.Close(); err != nil {
					// If the file cannot be closed, wait and try again.
					time.Sleep(100 * time.Millisecond)
				} else {
					// If the file is closed successfully, break the loop.
					break
				}
				if i == 4 {
					// If the file cannot be closed after 5 attempts, print an error message.
					fmt.Println("Failed to close file after 5 attempts")
				}
			}

			// Sync the file system to ensure data is written to disk.
			syscall.Sync()

			finishChan <- struct{}{}
		}()

		// Loop indefinitely to write data to the file.
		for {
			// Select on the data channel.
			select {
			case val, ok := <-dataChan:
				// If the channel is closed, exit the loop.
				if !ok {
					return
				}

				// Write each chunk of data to the pipe.
				for _, chunk := range val {
					// Write the chunk to the pipe.
					n, err := syscall.Write(pipe[1], chunk)
					if err != nil {
						// If the write fails, print an error message and exit.
						fmt.Printf("failed to write to pipe: %v\n", err)
						return
					}
					if n != len(chunk) {
						// If the write is partial, print an error message and exit.
						fmt.Printf("partial write to pipe, wrote %d bytes out of %d\n", n, len(chunk))
						return
					}

					// Splice the data from the pipe to the file.
					for n > 0 {
						written, err := syscall.Splice(pipe[0], nil, int(file.Fd()), nil, n, 0)
						if err != nil {
							// If the splice fails, print an error message and exit.
							fmt.Printf("failed to splice data: %v\n", err)
							return
						}
						n -= int(written)
					}
				}
			}
		}
	}()

	// Return the data channel and no error.
	return dataChan, finishChan, nil
}
