package utilhub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_RestoreAutoConfig demonstrates saving the default configuration to a JSON file.
// The second parameter (overwrite) is set to 'true', which means:
// - When true: The configuration will actually be written to the physical file
// - When false: The configuration is prepared but not written to disk
func Test_RestoreAutoConfig(t *testing.T) {
	// Setting it to true ensures the configuration is persisted to the filesystem.
	err := defaultConfig2file(&TestAutoConfig{}, false)
	require.NoError(t, err)
}
