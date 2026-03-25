package utilhub

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// =====================================================================================================================
// 🛠️ Linux Splice Tool
//
// Description:
// This tool provides a high-performance method for writing data to a file
// by leveraging the Linux splice() system call.
//
// Background:
// When handling large volumes of data, traditional read/write operations
// may introduce unnecessary memory copying between kernel space and user space.
// The splice() system call helps optimize this process by transferring data
// directly between file descriptors within the kernel.
//
// Purpose:
// - Efficiently write large amounts of data to files
// - Reduce CPU overhead and memory copying
// - Improve I/O performance in data-intensive scenarios
//
// Usage Scenario:
// Suitable for applications that require high-throughput data transfer,
// such as logging systems, data pipelines, or streaming services.
//
// =====================================================================================================================

// LinuxSpliceBulkWrite ⛏️ performs high-throughput bulk writes to a file by leveraging the Linux splice() system call.
// splice() is invoked to transfer data directly from the pipe to the destination file descriptor.
func LinuxSpliceBulkWrite(filename string, data [][]byte, fileFlag int, filePerm os.FileMode) error {
	// The file is opened with the specified flags and permissions.
	file, err := os.OpenFile(filename, fileFlag, filePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Create a pipe to send data through.
	pipe := make([]int, 2)
	if err = syscall.Pipe(pipe); err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	defer func() { _ = syscall.Close(pipe[0]) }()
	defer func() { _ = syscall.Close(pipe[1]) }()

	// Send each chunk of data through the pipe.
	for _, chunk := range data {
		var n int
		n, err = syscall.Write(pipe[1], chunk)
		if err != nil {
			return fmt.Errorf("failed to write to pipe: %w", err)
		}
		if n != len(chunk) {
			return fmt.Errorf("partial write to pipe, wrote %d bytes out of %d", n, len(chunk))
		}

		// The data is then written from the pipe to the file using the Splice system call.
		for n > 0 {
			var written int64
			written, err = syscall.Splice(pipe[0], nil, int(file.Fd()), nil, n, 0)
			if err != nil {
				return fmt.Errorf("failed to splice data: %w", err)
			}

			// Reduce n by the bytes written, so the loop continues until all data is written.
			// syscall.Splice 不保证一次就能写完所有资料，
			// n -= int(written) 的作用是「扣掉已经成功写入档案的资料，剩下的继续写」，
			// 如果没有这一行，n 永遠不會減小，迴圈就會變成 無限迴圈。
			n -= int(written)
		}
	}

	return nil
}

// LinuxSpliceBulkRead ⛏️ reads the contents of a file in chunks of the specified size.
// It returns a slice of byte slices, where each inner slice contains one chunk of data.
func LinuxSpliceBulkRead(filename string, chunkSize int) ([][]byte, error) {
	// Open the file for reading.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Ensure the file is closed when the function returns.
	defer func() { _ = file.Close() }()

	var chunks [][]byte
	// Allocate a buffer of the specified chunk size.
	buf := make([]byte, chunkSize)

	for {
		// Read up to chunkSize bytes from the file into the buffer.
		var n int
		n, err = file.Read(buf)
		if n == 0 {
			// End of file reached.
			break
		}
		if err != nil && err != io.EOF {
			// Return any read error other than EOF.
			return nil, err
		}

		// Copy the bytes read into a new slice to avoid overwriting the buffer.
		chunk := make([]byte, n)
		copy(chunk, buf[:n])

		// Append the chunk to the slice of chunks.
		chunks = append(chunks, chunk)
	}

	// Return all chunks read from the file.
	return chunks, nil
}

// GenerateBinaryFile ⛏️ creates a binary file at the specified path containing the given number of random UUIDs.
// Each UUID is 16 bytes, so the resulting file size will be count * 16 bytes.
func GenerateBinaryFile(path string, count int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	// Ensure the file is closed when the function returns.
	defer file.Close()

	for i := 0; i < count; i++ {
		// Generate a new UUID.
		id := uuid.New()
		// Write the 16-byte UUID to the file.
		_, err := file.Write(id[:])
		if err != nil {
			return err
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
	if err = syscall.Pipe(pipe); err != nil {
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
