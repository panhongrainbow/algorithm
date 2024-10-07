package utilhub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// auxiliaryTest is an auxiliary test function that takes an instance of type T.
var auxiliaryTest = func(tt T) {
	// Call the EnsureTestEnvironment method on the provided instance of T.
	tt.EnsureTestEnvironment()
}

// Test_TestOnlyFunc is used to test the EnsureTestEnvironment method.
func Test_TestOnlyFunc(t *testing.T) {
	// Call the auxiliary test function, passing a new instance of T.
	// This will execute the EnsureTestEnvironment method, which ensures that the function can only be called in a test environment.
	auxiliaryTest(T{})
}

// helperFunction is a helper function to test GetFunctionName.
func helperFunction(tt T) string {
	tt.EnsureTestEnvironment()
	return GetFunctionName()
}

// Test_GetFunctionName is A test function to call GetFunctionName.
// ⚠️ Please note that using an anonymous function will cause this method to fail. (使用匿名函式，会取值失败)
func Test_GetFunctionName(t *testing.T) {
	// Call the helper function and capture the result.
	result := helperFunction(T{})

	// Assert that the result is equal to the expected function name.
	expected := "utilhub.helperFunction"
	assert.Equal(t, expected, result, "Expected function name does not match")

	// Additional test to check calling GetFunctionName directly in the test function.
	result = GetFunctionName()
	expected = "utilhub.Test_GetFunctionName"
	assert.Equal(t, expected, result, "Expected function name does not match in the test function")

	// Testing when GetFunctionName is called from an anonymous function.
	anonFunc := func() string {
		return GetFunctionName()
	}
	result = anonFunc()
	expected = "utilhub.Test_GetFunctionName.func1"
	assert.Equal(t, expected, result, "Expected function name does not contain in the anonymous function")
}
