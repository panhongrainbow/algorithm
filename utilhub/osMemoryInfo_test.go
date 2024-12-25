package utilhub

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLinuxAvailableMemory(t *testing.T) {
	actualAvailableMemory, err := GetLinuxAvailableMemory()
	require.NoError(t, err)
	require.NotZero(t, actualAvailableMemory)

	// a := make([]int64, math.MaxInt64)
	// a[0] = 1
	// fmt.Println(a[0], len(a))

	// test1()
}
