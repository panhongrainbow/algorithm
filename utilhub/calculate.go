package utilhub

// =====================================================================================================================
//                  🛠️ Calculate (Tool)
// Calculate is a collection of functions for calculating various values.
// =====================================================================================================================

// Adjust2Even ⛏️ returns the specified length if it is even, or a slice of length+1 if it is odd.
// If the input length is negative, the function recursively calls itself with the absolute value of the length.
func Adjust2Even(length int64) int64 {
	// Check if the input length is a negative integer and
	if length < 0 {
		// Recursively call Adjust2Even with the absolute value of the length.
		length = Adjust2Even(-length)
		// If the length is negative, return the negative of the adjusted length.
		return -length
	}

	// Check if the length is even.
	if length%2 == 0 {
		// If the length is even, return the specified length.
		return length
	}

	// If the length is odd, return a slice of length+1 to maintain even.
	return length + 1
}
