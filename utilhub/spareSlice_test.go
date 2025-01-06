package utilhub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestCalculateAllocatableMemory tests the calculateAllocatableMemory function on a real system.
func Test_CalculateAllocatableMemory(t *testing.T) {
	// Define a slice of test cases, each containing a percentage and an expected error.
	tests := []struct {
		percentage  uint64 // The percentage of memory to allocate.
		expectedErr error  // The expected error for the test case.
	}{
		{50, nil},  // Test valid case with 50% memory allocation.
		{100, nil}, // Test valid case with 100% memory allocation.
		{0, nil},   // Test valid case with 0% (no allocation).
		{150, fmt.Errorf("invalid percentage: 150. Must be between 0 and 100")}, // Invalid percentage.
	}

	// Iterate over each test case.
	for _, tt := range tests {
		// Run a sub-test for each test case, with a descriptive name.
		t.Run(fmt.Sprintf("Percentage %d", tt.percentage), func(t *testing.T) {
			// Call the SpareSliceSize function with the test case's percentage.
			actualSize, actualErr := SpareSliceSize(tt.percentage)

			// Check if an error is expected for this test case.
			if tt.expectedErr != nil {
				// If an error is expected, verify that an error occurred.
				require.Error(t, actualErr)
				// Verify that the error message matches the expected error message.
				assert.EqualError(t, actualErr, tt.expectedErr.Error())
			} else {
				// If no error is expected, verify that no error occurred.
				require.NoError(t, actualErr)
			}

			// Verify that the allocated size is non-negative.
			assert.GreaterOrEqual(t, actualSize, uint64(0), "allocated size should be non-negative")
		})
	}
}
