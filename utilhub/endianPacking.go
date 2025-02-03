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
		if err := binary.Read(buf, order, &v); err != nil {
			// If an error occurs, break out of the loop as the reading is complete.
			break
		}

		// Append the successfully read int64 value to the result slice.
		slice = append(slice, v)
	}

	// Return the converted slice of int64 values and a nil error if the conversion was successful.
	return slice, nil
}
