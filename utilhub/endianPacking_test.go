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
