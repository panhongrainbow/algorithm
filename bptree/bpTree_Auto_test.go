package bpTree

// =====================================================================================================================
//                  ⚗️ Consistency Integrity Test ( [B Plus Tree] ) - B加树 主要测试
// =====================================================================================================================
// 🧪 B Plus Tree unit test validates structure via bulk insert and delete.
// 🧪 Inserts large data, then deletes all to check if tree resets to empty.
// 🧪 Indexing errors may cause data loss or deletion failures.
// 🧪 Ensures indexing accuracy for reliable operations.

// To run the test, run the following command:
//
// cd /home/panhong/go/src/github.com/panhongrainbow/go-algorithm/bptree
// go clean -cache
// go test -v . -timeout=0 -run Test_Check_BpTree_Accuracies

// =====================================================================================================================

import (
	"testing"

	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

// Test_Check_BpTree_Accuracy 🧫 checks if the tree resets after bulk insert/delete, ensuring indexing correctness.
func Test_Check_BpTree_Accuracies(t *testing.T) {
	// This is a set of path configurations shared by the automated testing.
	var (
		// 🧪 Create a config instance for B plus tree unit testing and parse default values.
		autoTestConfig = utilhub.GetAutoConfig()

		// 🧪 Navigate to the project dataSet directory for test record storage.
		ProjectDir = utilhub.FileNode{}.Goto(autoTestConfig.Record.TestRecordPath)

		// 🧪 Create a subdirectory named with the current date under the project.
		recordDir = ProjectDir.MkDir(autoTestConfig.Record.ManualRecordDate)
	)

	t.Run("Pre-test checks", func(t *testing.T) {
		// Record path must not be empty.
		require.NotEqual(t, "", ProjectDir.Path(), "record path is empty; check path creation")

		// Record subdirectory must not be empty.
		require.NotEqual(t, "", recordDir.Path(), "record date path is empty; check path creation")
	})

	t.Run("Bulk InsertDelete", func(t *testing.T) {

		// Test data will only be generated during automated testing.
		if autoTestConfig.Mechanism == "auto" {
			// Prepare test data for BulkInsertDelete.
			prepareBulkInsertDelete(t, recordDir, autoTestConfig)
		}

		// Only during automated testing or this test mode will the following tests be performed continuously.
		if autoTestConfig.Mechanism == "auto" {

			// Verify test data for BulkInsertDelete.
			verifyBulkInsertDelete(t, recordDir, autoTestConfig)

			// Execute accuracy test for BulkInsertDelete.
			runBulkInsertDelete(t, recordDir, autoTestConfig)
		}
	})

	t.Run("Randomized Boundary Test", func(t *testing.T) {

		// Test data will only be generated during automated testing.
		if autoTestConfig.Mechanism == "auto" {
			// Prepare test data for RandomizedBoundary.
			prepareRandomizedBoundary(t, recordDir)
		}

		// Only during automated testing or this test mode will the following tests be performed continuously.
		if autoTestConfig.Mechanism == "auto" {
			// Verify test data for RandomizedBoundary.
			verifyRandomizedBoundary(t, recordDir, autoTestConfig)

			// Execute accuracy test for RandomizedBoundary.
			runRandomizedBoundary(t, recordDir, autoTestConfig)
		}
	})

	t.Run("Single Node Endurance Test", func(t *testing.T) {

		// Test data will only be generated during automated testing.
		if autoTestConfig.Mechanism == "auto" {
			// Prepare test data for SingleNodeEndurance.
			prepareMode3(t, recordDir)
		}

		// Only during automated testing or this test mode will the following tests be performed continuously.
		if autoTestConfig.Mechanism == "auto" {
			// Verify test data for SingleNodeEndurance.
			verifySingleNodeEndurance(t, recordDir, autoTestConfig)

			// Execute accuracy test for SingleNodeEndurance.
			runSingleNodeEndurance(t, recordDir, autoTestConfig)
		}

	})
}
