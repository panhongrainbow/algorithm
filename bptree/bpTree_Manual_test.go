package bpTree

import (
	"fmt"
	"testing"

	"github.com/panhongrainbow/go-algorithm/utilhub"
)

func Test_Check_Manual_Accuracies(t *testing.T) {
	// This is a set of path configurations shared by the automated testing.
	var (
		// 🧪 Create a config instance for B plus tree unit testing and parse default values.
		manualTestConfig = utilhub.GetManualConfig()
	)

	// fmt.Println(manualTestConfig)

	for _, each := range manualTestConfig {
		var (
			// 🧪 Navigate to the project dataSet directory for test record storage.
			ProjectDir = utilhub.FileNode{}.Goto(each.Record.TestRecordPath)

			// 🧪 Create a subdirectory named with the date under the project.
			recordDir = ProjectDir.MkDir(each.Record.ManualRecordDate)
		)

		recordFile, err := recordDir.CheckFile(each.Record.ManualRecordFile)

		fmt.Println(recordFile.Path(), err)
	}

	return
}
