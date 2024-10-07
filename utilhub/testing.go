package utilhub

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ========================================================================================================================================
//                  🛠️ auxiliary Testing (Tool)
// The purpose of writing this function is to prevent the misuse of auxiliary functions by passing t *testing.T into them.
// However, using t *testing.T directly may introduce the locking mechanism inherent in testing.T, which could affect the test results.
// To address this, I created a custom type T struct{} and a method EnsureTestEnvironment to disguise it as t *testing.T.
// This design effectively prevents auxiliary functions from being misused in production code.
// ========================================================================================================================================

// T ⛏️ is a custom type T to disguise the testing type.
type T struct{}

// GetFunctionName ⛏️ returns the name of the caller's function, retaining only the last part.
func GetFunctionName() string {
	// Retrieve the program counter for the caller.
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		// Return "unknown" if the caller info is not available.
		return "unknown"
	}

	// Get the function information for the caller.
	function := runtime.FuncForPC(pc)
	if function == nil {
		// Return "unknown" if the function info is not available.
		return "unknown"
	}

	// Get the full function path.
	fullFunctionName := function.Name()

	// Find the last occurrence of '/' and the part after the last '.'
	slashIndex := strings.LastIndex(fullFunctionName, "/")
	if slashIndex == -1 {
		// If no slash is found, set to 0.
		slashIndex = 0
	} else {
		// Skip the last slash.
		slashIndex++
	}

	// Return the last part of the function name.
	return fullFunctionName[slashIndex:]
}

// EnsureTestEnvironment ⛏️ allows execution only in a test environment.
func (T) EnsureTestEnvironment() {
	// Get the executable's path and handle any errors.
	execName, err := os.Executable()
	if err != nil {
		panic("failed to get the executable path: " + err.Error())
	}

	// Resolve any symbolic links to get the actual path.
	resolvedPath, err := filepath.EvalSymlinks(execName)
	if err != nil {
		panic("failed to evaluate symlinks: " + err.Error())
	}

	// Check if the resolved path exists.
	// 很多函式都无法保证绝对路径，最后用 os.Stat 检查。
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		panic("the executable does not exist: " + resolvedPath)
	}

	// Check if the executable name contains any of the allowed test patterns.
	if !strings.Contains(resolvedPath, ".test") {
		// If none of the conditions match, trigger a panic.
		panic("this function can only be called in tests")
	}
}
