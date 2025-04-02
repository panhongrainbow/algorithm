package utilhub

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// =====================================================================================================================
//                  🛠️ EndianPacking (Tool)
// EndianPacking is a collection of functions for handling byte packing and unpacking with different endian formats.
// =====================================================================================================================

// >>>>> >>>>> One-Dimensional Functions ⛏️

// Int64SliceToBytes ⛏️ converts a slice of int64 values to a byte slice,
// order: The binary.ByteOrder to use for the conversion (e.g. binary.LittleEndian or binary.BigEndian).
func Int64SliceToBytes(slice []int64, order binary.ByteOrder) ([]byte, error) {
	// Check if the specified byte order is supported.
	if order != binary.LittleEndian && order != binary.BigEndian {
		return []byte{}, fmt.Errorf("unsupported byte order: %s", order)
	}

	// Handle the case where the input slice is empty or nil.
	if slice == nil || len(slice) == 0 {
		// Return an empty byte slice to indicate that no data was converted.
		return []byte{}, nil
	}

	// Create a new bytes buffer to store the converted data.
	var buf bytes.Buffer

	// Iterate over each int64 value in the input slice.
	for _, v := range slice {
		// Write the current int64 value to the buffer using the specified byte order.
		if err := binary.Write(&buf, order, v); err != nil {
			// If an error occurs during the write, return immediately with the error.
			return []byte{}, err
		}
	}

	// Return the converted byte slice and a nil error.
	return buf.Bytes(), nil
}

// BytesToInt64Slice ⛏️ converts a byte slice to a slice of int64 values based on the specified byte order.
// The function reads the input byte slice using the provided binary.ByteOrder (e.g., binary.LittleEndian or binary.BigEndian)
func BytesToInt64Slice(data []byte, order binary.ByteOrder) ([]int64, error) {
	// Check if the specified byte order is supported.
	if order != binary.LittleEndian && order != binary.BigEndian {
		return []int64{}, fmt.Errorf("unsupported byte order: %s", order)
	}

	// Handle the case where the input slice is empty or nil.
	if data == nil || len(data) == 0 {
		// Return an empty byte slice to indicate that no data was converted.
		return []int64{}, nil
	}

	// Check if the input data length is valid for conversion to int64 slice.
	// A valid length must be a multiple of 8 bytes (since int64 is 8 bytes long).
	if len(data) < 8 {
		// If the length is invalid, return an empty int64 slice and an error.
		return []int64{}, fmt.Errorf("invalid data length: %d (must be a multiple of 8)", len(data))
	}

	// Initialize an empty slice to store the converted int64 values.
	var slice []int64

	// Create a bytes reader from the input data to facilitate reading the byte slice.
	buf := bytes.NewReader(data)

	// Continuously read int64 values from the byte slice until an error occurs.
	for {
		// Declare a variable to store the current int64 value being read.
		var v int64

		// Attempt to read the next int64 value from the byte slice using the specified byte order.
		err := binary.Read(buf, order, &v)
		if err != nil && err.Error() != "EOF" {
			// If an error occurs, break out of the loop as the reading is complete.
			// break
			return []int64{}, err
		}

		if err != nil && err.Error() == "EOF" {
			break
		}

		// Append the successfully read int64 value to the result slice.
		slice = append(slice, v)
	}

	// Return the converted slice of int64 values and a nil error if the conversion was successful.
	return slice, nil
}

// >>>>> >>>>> Two-Dimensional Functions ⛏️

// Int64SliceToBlockBytes ⛏️ converts a byte slice to a 2 Dimensional slice of int64 values based on the specified byte order.
// The function reads the input byte slice using the provided binary.ByteOrder (e.g., binary.LittleEndian or binary.BigEndian)
func Int64SliceToBlockBytes(slice []int64, order binary.ByteOrder, startPoint, length, width int) (sliceBlock [][]byte, endPoint int, finished bool, err error) {
	// Check for invalid input parameters.
	if length <= 0 || width <= 0 {
		return nil, startPoint, false, fmt.Errorf("invalid length or width")
	}

	// Pre-allocate the sliceBlock to avoid append triggering reallocation.
	sliceBlock = make([][]byte, 0, length)

	// Iterate over the length to convert the int64 slice to byte blocks.
	for i := 0; i < length; i++ {
		// If the end point is equal to the length of the slice, there is no more data.
		if startPoint >= len(slice) {
			// Set the finished flag to true and break the loop.
			finished = true
			break
		}

		// Calculate the end point for the current block.
		var gotData []byte
		end := startPoint + width
		if end > len(slice) {
			// If the end point exceeds the length of the slice, set it to the length of the slice
			// and set the finished flag to true.
			end = len(slice)
			finished = true
		}

		// Convert the int64 slice to bytes using the provided byte order.
		gotData, err = Int64SliceToBytes(slice[startPoint:end], order)
		if err != nil {
			// Return an error if the conversion fails.
			return nil, startPoint, false, err
		}

		// Append the byte block to the sliceBlock.
		sliceBlock = append(sliceBlock, gotData)
		// Update the start point for the next block.
		startPoint = end
	}

	// Set the end point to the updated start point.
	endPoint = startPoint
	return
}
