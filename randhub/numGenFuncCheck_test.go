package randhub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_GenerateInt64Numbers tests the GenerateNumbers function with various scenarios.
func Test_GenerateInt64Numbers(t *testing.T) {
	// Define a list of test cases to cover different scenarios.
	tests := []struct {
		name      string // Name of the test case.
		count     uint64 // Number of integers to generate.
		minNum    int64  // Minimum value of the integers.
		maxNum    int64  // Maximum value of the integers.
		expectErr bool   // Whether an error is expected.
	}{
		{
			name:      "Valid range with sufficient numbers.",
			count:     10,
			minNum:    1,
			maxNum:    20,
			expectErr: false,
		},
		{
			name:      "Count exceeds range.",
			count:     20,
			minNum:    1,
			maxNum:    10,
			expectErr: false,
		},
		{
			name:      "minNum greater than maxNum.",
			count:     5,
			minNum:    10,
			maxNum:    1,
			expectErr: true,
		},
		{
			name:      "Single value range.",
			count:     1,
			minNum:    5,
			maxNum:    5,
			expectErr: false,
		},
		{
			name:      "Large range.",
			count:     1000,
			minNum:    1,
			maxNum:    10000,
			expectErr: false,
		},
	}

	// Iterate over each test case and run it.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the GenerateNumbers function with the current test case parameters.
			result, err := GenerateNumbersWithOption(tt.count, tt.minNum, tt.maxNum, false)

			if tt.expectErr {
				// If an error is expected, assert that an error was returned.
				assert.Error(t, err)
				// Also, assert that the result is nil in case of an error.
				assert.Nil(t, result)
			} else {
				// If no error is expected, assert that no error was returned.
				assert.NoError(t, err)
				// Assert that the length of the result matches the expected count.
				assert.Equal(t, tt.count, uint64(len(result)))

				// Check that all numbers in the result are within the specified range.
				for _, num := range result {
					assert.GreaterOrEqual(t, num, tt.minNum) // Assert that the number is greater than or equal to minNum.
					assert.LessOrEqual(t, num, tt.maxNum)    // Assert that the number is less than or equal to maxNum.
				}
			}
		})
	}
}

// Test_GenerateUniqueInt64Numbers tests the GenerateUniqueNumbers function with various scenarios.
func Test_GenerateUniqueInt64Numbers(t *testing.T) {
	// Define a list of test cases to cover different scenarios.
	tests := []struct {
		name      string // Name of the test case.
		count     uint64 // Number of unique integers to generate.
		minNum    int64  // Minimum value of the integers.
		maxNum    int64  // Maximum value of the integers.
		expectErr bool   // Whether an error is expected.
	}{
		{
			name:      "Valid range with sufficient numbers.",
			count:     10,
			minNum:    1,
			maxNum:    20,
			expectErr: false,
		},
		{
			name:      "Count exceeds range.",
			count:     20,
			minNum:    1,
			maxNum:    10,
			expectErr: true,
		},
		{
			name:      "minNum equals maxNum.",
			count:     1,
			minNum:    5,
			maxNum:    5,
			expectErr: false,
		},
		{
			name:      "minNum greater than maxNum.",
			count:     5,
			minNum:    10,
			maxNum:    1,
			expectErr: true,
		},
		{
			name:      "Single number requested from large range.",
			count:     1,
			minNum:    1,
			maxNum:    10000,
			expectErr: false,
		},
		{
			name:      "Large count within large range.",
			count:     1000,
			minNum:    1,
			maxNum:    10000,
			expectErr: false,
		},
	}

	// Iterate over each test case and run it.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the GenerateUniqueNumbers function with the current test case parameters.
			result, err := GenerateNumbersWithOption(tt.count, tt.minNum, tt.maxNum, true)

			if tt.expectErr {
				// If an error is expected, assert that an error was returned.
				assert.Error(t, err)
				// Also, assert that the result is nil in case of an error.
				assert.Nil(t, result)
			} else {
				// If no error is expected, assert that no error was returned.
				assert.NoError(t, err)
				// Assert that the length of the result matches the expected count.
				assert.Equal(t, tt.count, uint64(len(result)))

				// Create a map to track unique values in the result.
				uniqueNumbers := make(map[int64]struct{})

				// Check that all numbers in the result are unique and within the specified range.
				for _, num := range result {
					assert.GreaterOrEqual(t, num, tt.minNum) // Assert that the number is greater than or equal to minNum.
					assert.LessOrEqual(t, num, tt.maxNum)    // Assert that the number is less than or equal to maxNum.
					_, exists := uniqueNumbers[num]
					assert.False(t, exists, "Duplicate number found: %d", num) // Assert that the number is unique.
					uniqueNumbers[num] = struct{}{}
				}
			}
		})
	}
}

// Test_GenerateUniqueFloat64Numbers tests the GenerateUniqueNumbers function with various scenarios for float64 type.
func Test_GenerateUniqueFloat64Numbers(t *testing.T) {
	// Define a list of test cases to cover different scenarios.
	tests := []struct {
		name      string  // Name of the test case.
		count     uint64  // Number of unique float64 numbers to generate.
		minNum    float64 // Minimum value of the float64 numbers.
		maxNum    float64 // Maximum value of the float64 numbers.
		expectErr bool    // Whether an error is expected.
	}{
		{
			name:      "Valid range with sufficient numbers.",
			count:     10,
			minNum:    1.5,
			maxNum:    20.5,
			expectErr: false,
		},
		{
			name:      "Count exceeds range.",
			count:     20,
			minNum:    1.5,
			maxNum:    10.5,
			expectErr: false,
		},
		{
			name:      "minNum equals maxNum.",
			count:     1,
			minNum:    5.5,
			maxNum:    5.5,
			expectErr: false,
		},
		{
			name:      "minNum greater than maxNum.",
			count:     5,
			minNum:    10.5,
			maxNum:    1.5,
			expectErr: true,
		},
		{
			name:      "Single number requested from large range.",
			count:     1,
			minNum:    1.5,
			maxNum:    10000.75,
			expectErr: false,
		},
		{
			name:      "Large count within large range.",
			count:     1000,
			minNum:    1.25,
			maxNum:    10000.95,
			expectErr: false,
		},
	}

	// Iterate over each test case and run it.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the GenerateNumbersWithOption function with the current test case parameters.
			result, err := GenerateNumbersWithOption(tt.count, tt.minNum, tt.maxNum, true)

			if tt.expectErr {
				// If an error is expected, assert that an error was returned.
				assert.Error(t, err)
				// Also, assert that the result is nil in case of an error.
				assert.Nil(t, result)
			} else {
				// If no error is expected, assert that no error was returned.
				assert.NoError(t, err)
				// Assert that the length of the result matches the expected count.
				assert.Equal(t, int(tt.count), len(result))

				// Create a map to track unique values in the result.
				uniqueNumbers := make(map[float64]struct{})

				// Check that all numbers in the result are unique and within the specified range.
				for _, num := range result {
					assert.GreaterOrEqual(t, num, tt.minNum) // Assert that the number is greater than or equal to minNum.
					assert.LessOrEqual(t, num, tt.maxNum)    // Assert that the number is less than or equal to maxNum.
					_, exists := uniqueNumbers[num]
					assert.False(t, exists, "Duplicate number found: %f", num) // Assert that the number is unique.
					uniqueNumbers[num] = struct{}{}
				}
			}
		})
	}
}
