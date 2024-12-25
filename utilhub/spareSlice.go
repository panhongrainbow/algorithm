package utilhub

import (
	"fmt"
	"unsafe"
)

// =====================================================================================================================
//                  🛠️ Spare Slice (Tool)
// =====================================================================================================================

// memoryDecrementStep is the step size used to decrement the memory allocation size
// when attempting to allocate memory safely.
const memoryDecrementStep = 1024

// spareSlice ⛏️ calculates the maximum size of an array that can be allocated based on
// the specified percentage of available memory.
func spareSlice(percentage uint64) (uint64, error) {
	// Check if the specified percentage is valid (between 0 and 100).
	if percentage > 100 {
		// If the percentage is invalid, return an error.
		return 0, fmt.Errorf("invalid percentage: %d. Must be between 0 and 100", percentage)
	}

	// Get the available memory on the Linux system.
	availableMemory, err := GetLinuxAvailableMemory()
	if err != nil {
		// If an error occurs while getting the available memory, return the error.
		return 0, err
	}

	// Calculate the size of an int64 in bytes.
	int64SizeBytes := uint64(unsafe.Sizeof(int64(0)))

	// Calculate the maximum array size based on the available memory and specified percentage.
	maxArraySize := availableMemory * percentage / 100 * 1024 / int64SizeBytes

	// Attempt to allocate memory of the calculated size, decrementing the size by
	// memoryDecrementStep until a successful allocation is made.
	for ; maxArraySize >= 0; maxArraySize -= memoryDecrementStep {
		// Check if the memory allocation is successful.
		if allocationSuccessful := allocateMemorySafely(int(maxArraySize)); allocationSuccessful {
			// If the allocation is successful, break out of the loop.
			break
		}
	}

	// Return the maximum array size that was successfully allocated.
	return maxArraySize, nil
}

// allocateMemorySafely ⛏️ attempts to allocate memory of the specified size and returns
// true if successful, false otherwise.
func allocateMemorySafely(size int) bool {
	// Use a deferred function to catch any panics that occur during memory allocation.
	defer func() {
		if err := recover(); err != nil {
			// If a panic occurs, print the error message to the console.
			fmt.Print("Memory allocation failed: ", err)
		}
	}()

	// Attempt to allocate memory of the specified size.
	_ = make([]int64, size)

	// If no panic occurred, return true to indicate successful allocation.
	return true
}
