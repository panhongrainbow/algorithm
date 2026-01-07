package utilhub

import (
	"fmt"
	"testing"
)

func Test_ParseFailRecovery(t *testing.T) {
	err := ParseManual(&_manualTestConfig)
	test := _manualTestConfig
	fmt.Println(test)
	if err != nil {
		t.Error("fail recovery is not nil")
	}
}
