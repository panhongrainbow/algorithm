package utilhub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseToggle(t *testing.T) {
	// ParseToggle() is guarded by sync.Once.
	// Once the configuration has been initialized, the loading logic
	// cannot be executed again within the same test process.
	/* 不能这样测，因为被 sync.Once 组挡了
	var cfg BptreeToggleConfig
	err := ParseToggle(&cfg)
	require.NoError(t, err)
	require.True(t, cfg.Mechanism == "auto" || cfg.Mechanism == "manual")
	*/

	// Instead of calling ParseToggle() again, verify the value loaded during the first initialization.
	require.True(t, _toggleConfig.Mechanism == "auto" || _toggleConfig.Mechanism == "manual")
}
