package utilhub

import "testing"

func Test_ParseFailRecovery(t *testing.T) {
	err := ParseManual(&_manualTestConfig)
	if err != nil {
		t.Error("fail recovery is not nil")
	}
}
