package bpTree

import (
	"fmt"
	"testing"

	"github.com/panhongrainbow/go-algorithm/utilhub"
)

// This is a set of path configurations shared by the manual testing.
var (
	// 🧪 Create a config instance for B plus tree unit testing and parse default values.
	manualTestConfig = utilhub.GetManualConfig()

	// 🧪 Navigate to the project dataSet directory for test record storage.
	// ProjectDir = utilhub.FileNode{}.Goto(autoTestConfig.Record.TestRecordPath)

	// 🧪 Create a subdirectory named with the current date under the project.
	// recordDir4 = ProjectDir.MkDir(autoTestConfig.Record.ManualRecordDate)
)

func Test_Check_Manual_Accuracies(t *testing.T) {
	fmt.Printf("%+v\n", utilhub.GetManualConfig())
}
