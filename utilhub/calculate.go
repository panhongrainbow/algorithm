package utilhub

// =====================================================================================================================
//                  🛠️ Calculate (Tool)
// Calculate is a collection of functions for calculating various values.
// =====================================================================================================================

// Adjust2Even ⛏️ returns the specified length if it is even, or a slice of length+1 if it is odd.
// If the input length is negative, the function recursively calls itself with the absolute value of the length.
/*
When the test scenario requires data to be paired, grouped in twos, or evenly distributed,
the test data should typically consist of an even number of records.
當測試资料必须成对、两两分组或平均分配时使用
*/
func Adjust2Even(length int64) int64 {
	// Check if the input length is a negative integer.
	if length < 0 {
		// Recursively call Adjust2Even with the absolute value of the length.
		length = Adjust2Even(-length) // 这两行在测试时，会把 -5 的输入回传成 -6 的输出。
		// If the length is negative, return the negative of the adjusted length.
		return -length // 这两行在测试时，会把 -5 的输入回传成 -6 的输出。
	}

	// Check if the length is even.
	if length%2 == 0 {
		// If the length is even, return the specified length.
		return length
	}

	// If the length is odd, return a slice of length+1 to maintain even.
	return length + 1
}
