package bpTree

import (
	"strings"
	"testing"

	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/require"
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

		// Verify that the record file exists.
		if _, err := recordDir.CheckFile(each.Record.ManualRecordFile); err != nil {
			require.NoError(t, err)
		}

		// fmt.Println(recordFile.Path(), err)

		// Basic test.
		if strings.Contains(each.Record.ManualRecordFile, "BulkInsertDelete") {
			verifyBulkInsertDelete(t, recordDir, each)
		}

		// Boundary test.
		if strings.Contains(each.Record.ManualRecordFile, "RandomizedBoundary") {
			verifyRandomizedBoundary(t, recordDir, each)
		}

		// Endurance test.
		if strings.Contains(each.Record.ManualRecordFile, "SingleNodeEndurance") {
			verifySingleNodeEndurance(t, recordDir, each)
		}

	}

	return
}
