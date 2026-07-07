package utilhub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_RestoreAutoConfig demonstrates saving the default configuration to a JSON file.
// The second parameter (overwrite) is set to 'true', which means:
// - When true: The configuration will actually be written to the physical file.
// - When false: The configuration is prepared but not written to disk. (会把设定覆盖自动设定档！)
func Test_RestoreAutoConfig(t *testing.T) {
	// Setting it to true ensures the configuration is persisted to the filesystem.
	err := autoConfig2file(&AutoConfigType{}, false)

	// Check whether the generated default values are incorrect.
	require.NoError(t, err)
}

// Preview the current auto-configuration values for inspection.
func Test_PreviewAutoConfig(t *testing.T) {
	fmt.Printf("%+v\n", _autoConfig)
}
