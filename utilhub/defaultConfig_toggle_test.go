package utilhub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseToggle(t *testing.T) {
	var cfg ToggleConfig
	err := ParseToggle(&cfg)
	require.NoError(t, err)
	fmt.Printf("%+v\n", cfg)
}
