package utilhub

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// =====================================================================================================================
//                  🛠️ SpliceProgressNode (Tool)
// [SpliceProgressNode] are tools that combine [SpliceNode Package] and [ProgressBar Package].
// =====================================================================================================================

// LinuxSpliceProgressStreamWrite is a function that writes data to a file using Linux splicing and displays a progress bar.
func (fn FileNode) LinuxSpliceProgressStreamWrite(
	barTitle, barColor string, // <----- for ProgressBar function.
	filename string, fileFlag int, filePerm os.FileMode, // <----- for LinuxSpliceStreamWrite function.
	testDataSet []int64, // <----- for Int64SliceToBlockBytes function.
	order binary.ByteOrder, spliceBlockLength, spliceBlockWidth int, // <----- for Int64SliceToBlockBytes function.
) error {
	// #################################################################################################
	// Initialize linux splice stream writer and set up some parameters for data writing.
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
	// Initialize the progress bar.
	// #################################################################################################

	// Create a progress bar with optional configurations.
	progressBar, err := NewProgressBar(
		barTitle,                    // Progress bar title.
		uint32(len(testDataSet)),    // Total number of operations.
		70,                          // Progress bar width.
		WithTracking(5),             // Update interval.
		WithTimeZone("Asia/Taipei"), // Time zone.
		WithTimeControl(500),        // Update interval in milliseconds.
		WithDisplay(barColor),       // Display style.
	)

	if err != nil {
		return fmt.Errorf("failed to create progress bar: %w", err)
	}

	// Start the progress bar printer in a separate goroutine.
	go func() {
		progressBar.ListenPrinter()
	}()

	// Write data to the file in blocks until the entire data set is written.
	for !spliceWritingFinished {

		// #################################################################################################
		// Convert the data set to a block of bytes.
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

		time.Sleep(1 * time.Second)

		// Update the progress bar with the number of bytes written.
		progressBar.AddSpecificTimes(uint32(spliceBlockLength * spliceBlockWidth))
	}

	// #################################################################################################
	// Wait for the finish channel to receive the finish signal.
	// #################################################################################################

	// Mark the progress bar as complete.
	progressBar.Complete()

	// Wait for the progress bar printer to stop.
	<-progressBar.WaitForPrinterStop()

	// Close the data channel to signal the end of writing.
	close(spliceDataChan)

	// Wait for the finish channel to receive the finish signal.
	<-spliceFinishChan

	return nil
}
