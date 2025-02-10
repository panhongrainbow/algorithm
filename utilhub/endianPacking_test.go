package utilhub

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestInt64SliceToBytes tests the Int64SliceToBytes function, which converts a slice of int64 values to a byte slice.
func Test_Int64SliceToBytes(t *testing.T) {
	// Define a slice of test cases, each containing the input parameters and expected results.
	tests := []struct {
		name     string
		slice    []int64
		order    binary.ByteOrder
		wantErr  bool
		wantData []byte
	}{
		{
			// Test case: Little Endian byte order.
			name:     "Little Endian",
			slice:    []int64{1, 2, 3},
			order:    binary.LittleEndian,
			wantErr:  false,
			wantData: []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			// Test case: Big Endian byte order.
			name:     "Big Endian",
			slice:    []int64{1, 2, 3},
			order:    binary.BigEndian,
			wantErr:  false,
			wantData: []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3},
		},
		{
			// Test case: Empty input slice.
			name:     "Empty Slice",
			slice:    []int64{},
			order:    binary.LittleEndian,
			wantErr:  false,
			wantData: []byte{},
		},
		{
			// Test case: Invalid byte order (nil).
			name:     "Error Order",
			slice:    []int64{},
			order:    nil,
			wantErr:  true,
			wantData: []byte{},
		},
	}

	// Iterate over each test case and run a sub-test.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the Int64SliceToBytes function with the input parameters.
			gotData, err := Int64SliceToBytes(tt.slice, tt.order)

			// Check if an error is expected.
			if tt.wantErr {
				// If an error is expected, verify that one was returned.
				require.Error(t, err, "Expected error but got none")
			} else {
				// If no error is expected, verify that none was returned.
				require.NoError(t, err, "Unexpected error: %v", err)
			}

			// Verify that the returned byte slice matches the expected result.
			assert.Equal(t, tt.wantData, gotData, "Mismatch in byte conversion")
		})
	}
}

// Test_BytesToInt64Slice tests the BytesToInt64Slice function, which converts a byte slice to a slice of int64 values based on the specified byte order.
func Test_BytesToInt64Slice(t *testing.T) {
	// Define a slice of test cases, each containing the input parameters and expected results.
	tests := []struct {
		name    string
		data    []byte
		order   binary.ByteOrder
		want    []int64
		wantErr bool
	}{
		{
			// Test case: Little Endian byte order.
			name:  "Little Endian",
			data:  []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0},
			order: binary.LittleEndian,
			want:  []int64{1, 2},
		},
		{
			// Test case: Big Endian byte order.
			name:  "Big Endian",
			data:  []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2},
			order: binary.BigEndian,
			want:  []int64{1, 2},
		},
		{
			// Test case: Invalid data (length not a multiple of 8).
			name:    "Invalid Data",
			data:    []byte{1, 2, 3},
			order:   binary.LittleEndian,
			want:    []int64{},
			wantErr: true,
		},
		{
			// Test case: Empty data.
			name:  "Empty Data",
			data:  []byte{},
			order: binary.LittleEndian,
			want:  []int64{},
		},
	}

	// Iterate over each test case and run a sub-test.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the BytesToInt64Slice function with the input parameters.
			got, err := BytesToInt64Slice(tt.data, tt.order)

			// Check if an error is expected.
			if tt.wantErr {
				// If an error is expected, assert that one was returned.
				require.Error(t, err, "Expected error but got none")
			} else {
				// If no error is expected, ensure the function completed successfully.
				require.NoError(t, err, "Unexpected error: %v", err)
				// Assert that the result matches the expected value.
				assert.Equal(t, tt.want, got, "Mismatch in decoded int64 slice")
			}
		})
	}
}

// Test_Int64SliceToBlockBytes tests the Int64SliceToBlockBytes function to ensure it correctly converts slices of int64 into blocks of bytes.
func Test_Int64SliceToBlockBytes(t *testing.T) {
	tests := []struct {
		name       string           // Test case name
		slice      []int64          // Input slice of int64 numbers
		order      binary.ByteOrder // Byte order (LittleEndian or BigEndian)
		startPoint int              // The starting index in the slice
		length     int              // Number of blocks to be generated
		width      int              // Number of int64 values per block
		want       [][]byte         // Expected output blocks
		wantErr    bool             // Whether an error is expected
	}{
		{
			name:       "LittleEndian Block",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.LittleEndian,
			startPoint: 0,
			length:     2,
			width:      4,
			want: [][]byte{
				{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0},
				{5, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0},
			},
			wantErr: false,
		},
		{
			name:       "BigEndian Block",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.BigEndian,
			startPoint: 0,
			length:     2,
			width:      4,
			want: [][]byte{
				{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4},
				{0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 8},
			},
			wantErr: false,
		},
		{
			name:       "Partial Block - BigEndian",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.BigEndian,
			startPoint: 5,
			length:     2,
			width:      4,
			want: [][]byte{
				{0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 8},
			},
			wantErr: false,
		},
		{
			name:       "Invalid Length",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.BigEndian,
			startPoint: 0,
			length:     0, // Invalid: length cannot be zero
			width:      4,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "Invalid Width",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.BigEndian,
			startPoint: 0,
			length:     2,
			width:      0, // Invalid: width cannot be zero
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "Start Point Exceeds Slice Length",
			slice:      []int64{1, 2, 3, 4, 5, 6, 7, 8},
			order:      binary.BigEndian,
			startPoint: 10, // Invalid: startPoint is beyond slice length
			length:     2,
			width:      4,
			want:       [][]byte{}, // Expect an empty result
			wantErr:    false,
		},
	}

	// Iterate over all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			got, _, _, err := Int64SliceToBlockBytes(
				tt.slice, tt.order, tt.startPoint, tt.length, tt.width,
			)

			// Check if an error was expected
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got) // Expect nil output on error
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got) // Verify output matches expectation
			}
		})
	}
}
