package utilhub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_Adjust2Even tests the Even function with various input lengths.
func Test_Adjust2Even(t *testing.T) {
	// Define a slice of test cases, each containing a name, input length, and expected output.
	tests := []struct {
		name     string // Name of the test case.
		length   int64  // Input length to be passed to the Even function.
		expected int64  // Expected output of the Even function.
	}{
		// Test case for an even length.
		{"even length", 6, 6},
		// Test case for an odd length.
		{"odd length", 7, 8},
		// Test case for a negative length.
		{"negative length", -5, -6},
		// Test case for a zero length.
		{"zero length", 0, 0},
	}

	// Iterate over each test case and run a sub-test.
	for _, test := range tests {
		// Run a sub-test with the test case name.
		t.Run(test.name, func(t *testing.T) {
			// Call the Even function with the test case input length.
			actual := Adjust2Even(test.length)
			// Assert that the actual output matches the expected output.
			assert.Equal(t, test.expected, actual)
		})
	}
}
