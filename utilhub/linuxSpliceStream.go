package utilhub

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// =====================================================================================================================
// 🛠️ Linux Splice Stream Tool
//
// Description:
// This tool implements a high-performance streaming data transfer mechanism
// using the Linux splice() system call, enabling zero-copy data movement
// between file descriptors.
//
// Background:
// Traditional read/write operations require copying data between kernel space
// and user space, which introduces CPU overhead and memory bandwidth usage.
// The splice() system call avoids these extra copies by transferring data
// directly within the kernel via pipes.
//
// Purpose:
// - Stream large volumes of data efficiently
// - Minimize CPU usage and memory copying
// - Enable high-throughput, low-latency data pipelines
//
// Usage Scenario:
// Ideal for:
// - File-to-file streaming
// - Socket-to-file (e.g., download writers)
// - File-to-socket (e.g., static file servers)
// - Logging pipelines and real-time data forwarding
//
// =====================================================================================================================

// LinuxSpliceStreamWrite ⛏️ creates a pipe to write data to a file using the Splice system call.
// It returns a channel to send data to be written to the file.
func LinuxSpliceStreamWrite(filename string, fileFlag int, filePerm os.FileMode) (dataChan chan [][]byte, finishChan chan struct{}, err error) {
	// Open the file with the specified flags and permissions.
	var file *os.File
	file, err = os.OpenFile(filename, fileFlag, filePerm)
	if err != nil {
		// If the file cannot be opened, return an error.
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Create a pipe to write data to.
	pipe := make([]int, 2)
	if err = syscall.Pipe(pipe); err != nil {
		// If the pipe cannot be created, close the file and return an error.
		_ = file.Close()
		return nil, nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	// Create a channel to send data to be written to the file.
	dataChan = make(chan [][]byte, 100)

	// Signal that a process has finished.
	finishChan = make(chan struct{})

	// Start a goroutine to write data to the file.
	go func() {
		// Defer closing the pipe and file.
		defer func() {
			// Close the pipe.
			_ = syscall.Close(pipe[0])
			_ = syscall.Close(pipe[1])

			// Close the file with retries.
			for i := 0; i < 20; i++ {
				if closeErr := file.Close(); closeErr != nil {
					// If the file cannot be closed, wait and try again.
					time.Sleep(1 * time.Second)
				} else {
					// If the file is closed successfully, break the loop.
					break
				}
				if i == 19 {
					// If the file cannot be closed after 5 attempts, print an error message.
					fmt.Println("Failed to close file after 20 attempts")
				}
			}

			// Sync the file system to ensure data is written to disk.
			syscall.Sync()

			// Send a signal to indicate that the process has finished.
			finishChan <- struct{}{}
		}()

		// Loop indefinitely to write data to the file.
		for {
			// Select on the data channel.
			select {

			// Receive a value from dataChan and check whether the channel is still open.
			case val, ok := <-dataChan:

				// If the channel is closed, exit the loop.
				if !ok {
					return
				}

				// Write each chunk of data to the pipe.
				for _, chunk := range val {

					// Write the chunk to the pipe.
					var n int
					n, err = syscall.Write(pipe[1], chunk)
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
						var written int64
						written, err = syscall.Splice(pipe[0], nil, int(file.Fd()), nil, n, 0)
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

// LinuxSpliceStreamRead ⛏️ efficiently reads a file in chunks using syscall.Splice, splitting the data into a 2D byte slice and sending it through a channel for further processing.
// It returns a channel that delivers data read from the file in chunks.
func LinuxSpliceStreamRead(filename string, fileFlag int, filePerm os.FileMode, chunkSize int, chunkRaw int) (dataChan chan [][]byte, finishChan chan struct{}, err error) {
	// Open the file with the specified flags and permissions.
	var file *os.File
	file, err = os.OpenFile(filename, fileFlag, filePerm)
	if err != nil {
		// If the file cannot be opened, return an error.
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Create a pipe to write data to.
	pipe := make([]int, 2)
	if err = syscall.Pipe(pipe); err != nil {
		// If the pipe cannot be created, close the file and return an error.
		_ = file.Close()
		return nil, nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	// Create a channel to send data to be read from the file.
	dataChan = make(chan [][]byte, 100)

	// Signal that a process has finished.
	finishChan = make(chan struct{})

	// Start a goroutine to write data to the file.
	go func() {
		// Defer closing the pipe and file.
		defer func() {
			// Close the pipe.
			_ = syscall.Close(pipe[0])
			_ = syscall.Close(pipe[1])
		}()

		// Close the file with retries.
		for i := 0; i < 20; i++ {
			if closeErr := file.Close(); closeErr != nil {
				// If the file cannot be closed, wait and try again.
				time.Sleep(1 * time.Second)
			} else {
				// If the file is closed successfully, break the loop.
				break
			}
			if i == 19 {
				// If the file cannot be closed after 5 attempts, print an error message.
				fmt.Println("Failed to close file after 20 attempts")
			}
		}

		// Send a signal to indicate that the process has finished.
		finishChan <- struct{}{}
	}()

	// Accumulate data.
	var accumulatedData []byte

	// Start reading.
	for {
		// Use splice to transfer data from file to pipe.
		n, err := syscall.Splice(int(file.Fd()), nil, pipe[1], nil, chunkSize*chunkRaw, 0)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		// Read data from pipe.
		buf := make([]byte, n)
		total := 0
		for total < int(n) {
			m, err := unix.Read(pipe[0], buf[total:])
			if err != nil {
				return nil, nil, err
			}
			if m == 0 {
				break
			}
			total += m
		}

		// Accumulate the data.
		if total > 0 {
			accumulatedData = append(accumulatedData, buf[:total]...)
		}

		// If enough data is accumulated, split it into chunks of chunkSize by chunkSize2.
		for len(accumulatedData) >= chunkSize*chunkRaw {
			// Take chunkSize * chunkSize2 bytes.
			block := accumulatedData[:chunkSize*chunkRaw]
			accumulatedData = accumulatedData[chunkSize*chunkRaw:] // Remove processed data.

			// Split the block into chunkSize rows, each with chunkSize2 bytes.
			var chunk [][]byte
			for i := 0; i < chunkSize; i++ {
				start := i * chunkRaw
				end := start + chunkRaw
				chunk = append(chunk, block[start:end])
			}

			// Send the chunk.
			dataChan <- chunk
		}
	}

	// If there is remaining data, send it out.
	if len(accumulatedData) > 0 {
		chunk := [][]byte{accumulatedData}
		dataChan <- chunk
	}

	return dataChan, finishChan, nil
}
