package utilhub

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
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

// LinuxSpliceBulkRead ⛏️ reads the contents of a file using Linux splice.
// It returns a slice of byte slices, where each inner slice contains one chunk of data.
func LinuxSpliceBulkRead(filename string, chunkSize int) ([][]byte, error) {
	// Open the file for reading.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Ensure the file is closed when the function returns.
	defer func() { _ = file.Close() }()

	// Create a pipe to send data through.
	pipe := make([]int, 2)
	if err = syscall.Pipe(pipe); err != nil {
		return [][]byte{}, fmt.Errorf("failed to create pipe: %w", err)
	}
	defer func() { _ = syscall.Close(pipe[0]) }()
	defer func() { _ = syscall.Close(pipe[1]) }()

	// Store all the chunks read from the file.
	var chunks [][]byte

	for {
		// Use splice to transfer data from the file to the pipe.
		var n int64
		n, err = syscall.Splice(int(file.Fd()), nil, pipe[1], nil, chunkSize, 0)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Allocate a buffer to hold the chunk in user space
		buf := make([]byte, n)
		total := 0

		// Loop until we read all n bytes from the pipe
		for int64(total) < n {
			var m int
			m, err = unix.Read(pipe[0], buf[total:])
			if err != nil {
				return nil, err
			}
			if m == 0 {
				break
			}
			total += m
		}

		// Append the fully-read buffer slice to the chunks
		if total > 0 {
			chunks = append(chunks, buf[:total])
		}
	}

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
		_, err = file.Write(id[:])
		if err != nil {
			return err
		}
	}

	return nil
}
