package utilhub

import (
	"fmt"
	"testing"
)

// Preview the current manual configuration values for inspection.
func Test_PreviewManualConfig(t *testing.T) {
	fmt.Printf("%+v\n", _manualTestConfig)
}
