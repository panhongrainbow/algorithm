package utilhub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestForceReloadConfig verifies that the configuration is restored after being modified and reloaded.
func TestForceReloadConfig(t *testing.T) {
	// Reload configuration and initialize the test configuration state.
	ForceReloadAutoTest()

	// Modify the configuration value to simulate a changed runtime state.
	modifiedValue := GetAutoRandomTotalCount() + 1

	SetAutoRandomTotalCount(modifiedValue)

	// Reload configuration and verify that the modified value is not persisted.
	ForceReloadAutoTest()

	// Verify that the modified value is not retained after configuration reload.
	require.NotEqual(t, modifiedValue, GetAutoRandomTotalCount())
}
