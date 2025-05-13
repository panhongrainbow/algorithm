package utilhub

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// =====================================================================================================================
//                  🛠️ SpliceProgressNode (Tool)
// [SpliceProgressNode] are tools that combine [SpliceNode Package] and [ProgressBar Package]. (混合工具，锚钉文件 和 进度条 合拼)
// These functions are classified into two types, one for writing data to a file using Linux splicing and another for reading data from a file by chunks.
// LinuxSpliceProgressStreamWrite is for writing data to a file using Linux splicing and displaying a progress bar. (这是写入部份)
// ReadAllBytesWithProgress reads the entire content of a file into memory and displays a progress bar indicating the read progress.
// ReadBytesInChunksWithProgress reads data from a file in chunks, displaying a progress bar as it processes each chunk.
// =====================================================================================================================

// LinuxSpliceProgressStreamWrite is a function that writes data to a file using Linux splicing and displays a progress bar.
func (fn FileNode) LinuxSpliceProgressStreamWrite(
	// [Inputs]
	// <----- original data
	testDataSet []int64, // 原始资料
	// <----- parameters for writing
	filename string, fileFlag int, filePerm os.FileMode, // 写入资料要用到的参数
	order binary.ByteOrder, spliceBlockLength, spliceBlockWidth int,
	// <----- for ProgressBar function
	barTitle, barColor string, barLength int, // 进度条参数
) error { // [Outputs]

	// #################################################################################################
	// Initialize linux splice stream writer and set up some parameters for data writing. (初始化)
	// #################################################################################################

	// Initialize a Linux splice stream writer to write data to a file.
	// The file is created with write-only permissions and truncated if it already exists.
	spliceDataChan, spliceFinishChan, err := fn.LinuxSpliceStreamWrite(filename, fileFlag, filePerm)
	if err != nil {
		return fmt.Errorf("failed to initialize linux splice stream writer: %w", err)
	}

	// Variable Parameters:
	var (
		spliceWritingPoint    = 0     // Initialize the start point for block writing.
		spliceWritingFinished = false // Initialize a flag to track whether the writing process is finished.
	)

	// #################################################################################################
	// Initialize the progress bar. (准备进度条)
	// #################################################################################################

	// ▓▒░ Create a progress bar with optional configurations.
	progressBar, err := NewProgressBar(
		barTitle,                    // Progress bar title.
		uint32(len(testDataSet)),    // Total number of operations.
		barLength,                   // Progress bar width.
		WithTracking(5),             // Update interval.
		WithTimeZone("Asia/Taipei"), // Time zone.
		WithTimeControl(500),        // Update interval in milliseconds.
		WithDisplay(barColor),       // Display style.
	)

	if err != nil {
		return fmt.Errorf("failed to create progress bar: %w", err)
	}

	// ▓▒░ Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// Write data to the file in blocks until the entire data set is written.
	for !spliceWritingFinished {

		// #################################################################################################
		// Convert the data set to a block of bytes. (把数据转换为字节块，并决定端序)
		// #################################################################################################

		// Convert the data set to a block of bytes using the Int64SliceToBlockBytes method in utilhub.
		// This method converts a slice of int64 values to a block of bytes.
		var block [][]byte
		block, spliceWritingPoint, spliceWritingFinished, err = Int64SliceToBlockBytes(testDataSet, order, spliceWritingPoint, spliceBlockLength, spliceBlockWidth)
		// Check if an error occurred during block writing.
		if err != nil {
			return fmt.Errorf("failed to convert data to block: %w", err)
		}
		// Write the block to the file using the data channel.
		spliceDataChan <- block

		// Update the progress bar with the number of bytes written.
		progressBar.AddSpecificTimes(uint32(spliceBlockLength * spliceBlockWidth))
	}

	// #################################################################################################
	// 1. Wait for [the finish channel] to receive [the finish signal]. (等待写入完成)
	// 2. Wait for [the progress bar] to finish. (等待进度条完成)
	// #################################################################################################

	// -----> for the progress bar.

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// -----> for the finish channel.

	//Close the data channel to signal the end of writing.
	close(spliceDataChan)

	// Wait for the finish channel to receive the finish signal.
	<-spliceFinishChan

	// Return nil if the writing process is successful.
	return nil
}

// ReadAllBytesWithProgress is a function that reads the entire content of a file into memory and displays a progress bar indicating the read progress.
func (fn FileNode) ReadAllBytesWithProgress(
	// [Inputs]
	// <----- original data
	dataLength uint32, // 资料长度
	// <----- parameters for reading
	filename string, // 档名
	chunkSize int, // 快取大小
	order binary.ByteOrder, // 端序
	// <----- parameters for progress bar
	barTitle, barColor string, barLength int,
) (
	// [Outputs]
	testDataSet []int64,
	err error,
) {
	// #################################################################################################
	// Initialize the progress bar. (准备进度条)
	// #################################################################################################

	// Create a progress bar with optional configurations.
	progressBar, err := NewProgressBar(
		barTitle,                    // Progress bar title.
		dataLength,                  // Total number of operations.
		barLength,                   // Progress bar width.
		WithTracking(5),             // Update interval.
		WithTimeZone("Asia/Taipei"), // Time zone.
		WithTimeControl(500),        // Update interval in milliseconds.
		WithDisplay(barColor),       // Display style.
	)

	if err != nil {
		return []int64{}, fmt.Errorf("failed to create progress bar: %w", err)
	}

	// Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// #################################################################################################
	// Read data from the file in chunks until the entire data set is read. (开始读取数据)
	// #################################################################################################

	// Read data from the file in chunks and update the progress bar accordingly.
	var result []int64
	dataChan, errChan := fn.ReadBytesInChunks(filename, chunkSize)

	// Continuously read data from the file until the entire data set is read.
Loop:
	for {
		// Select from the data and error channels to handle incoming data or errors.
		select {
		case err := <-errChan:
			// If an EOF error is received, break out of the loop to indicate the end of the data set.
			if err == io.EOF {
				break Loop
			}
			// If a non-EOF error occurs, return an error to indicate an unexpected issue during reading.
			if err != nil && err != io.EOF {
				return []int64{}, fmt.Errorf("unexpected error while reading: %w", err)
			}
		case rawData := <-dataChan:
			// Convert the raw data to a slice of int64 values using the provided byte order.
			data, _ := BytesToInt64Slice(rawData, order)

			// Append the converted data to the result slice.
			result = append(result, data...)

			// Update the progress bar with the number of bytes written.
			progressBar.AddSpecificTimes(uint32(len(result)))
		}
	}

	// #################################################################################################
	// Wait for the progress bar to finish. (等待进度条完成)
	// #################################################################################################

	// ▓▒░ Mark the progress bar as complete.
	progressBar.Complete()

	// ▓▒░ Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Return the result if there is no error.
	return result, nil
}

// ReadBytesInChunksWithProgress is a function that reads data from a file by chunks.
func (fn FileNode) ReadBytesInChunksWithProgress(
	// [Inputs]
	// <----- parameters for reading
	filename string, // 档名
	chunkSize int, // 快取大小
	order binary.ByteOrder, // 端序
) (
	// [Outputs]
	output chan []int64,
	errOutput chan error,
	finishChan chan struct{},
) {
	output = make(chan []int64, chunkSize/8)
	errOutput = make(chan error)
	finishChan = make(chan struct{})

	dataChan, errChan := fn.ReadBytesInChunks(filename, chunkSize)

	// Continuously read data from the file until the entire data set is read.
	go func() {
	Loop:
		for {
			// Select from the data and error channels to handle incoming data or errors.
			select {
			case err := <-errChan:
				// If an EOF error is received, break out of the loop to indicate the end of the data set.
				if err == io.EOF {
					finishChan <- struct{}{}
					break Loop
				}
				// If a non-EOF error occurs, return an error to indicate an unexpected issue during reading.
				if err != nil && err != io.EOF {
					errOutput <- fmt.Errorf("unexpected error while reading: %w", err)
				}
			case rawData := <-dataChan:
				// Convert the raw data to a slice of int64 values using the provided byte order.
				data, _ := BytesToInt64Slice(rawData, order)

				// Append the converted data to the result slice.
				output <- data
			}
		}
	}()

	// Return the result if there is no error.
	return
}
