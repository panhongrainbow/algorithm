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
	"math/rand"
	"testing"
	"time"

	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

var (
	// 🧪 Create a config instance for B plus tree unit testing and parse default values.
	unitTestConfig = utilhub.GetDefaultConfig()

	// 🧪 Navigate to the project dataSet directory for test record storage.
	ProjectDir = utilhub.FileNode{}.Goto(unitTestConfig.Record.TestRecordPath)

	// 🧪 Create a subdirectory named with the current date under the project.
	recordDir = ProjectDir.MkDir(time.Now().Format("2006-01-02"))
)

// Test_Check_BpTree_Accuracy 🧫 checks if the tree resets after bulk insert/delete, ensuring indexing correctness.
func Test_Check_BpTree_Accuracies(t *testing.T) {
	t.Run("Pre-test checks", func(t *testing.T) {
		// Record path must not be empty.
		require.NotEqual(t, "", ProjectDir.Path(), "record path is empty; check path creation")

		// Record subdirectory must not be empty.
		require.NotEqual(t, "", recordDir.Path(), "record date path is empty; check path creation")
	})

	/*
		t.Run("Mode 3:  Test", func(t *testing.T) {
			// Prepare test data for mode 3.
			// prepareMode3(t)

			// Verify test data for mode 3.
			// verifyMode3(t)

			// Execute accuracy test for mode 3.
			runMode3(t)
		})

		return
	*/

	t.Run("Mode 1: Bulk Insert/Delete", func(t *testing.T) {
		// Prepare test data for mode 1.
		prepareMode1(t)

		// Verify test data for mode 1.
		verifyMode1(t)

		// Execute accuracy test for mode 1.
		runMode1(t)
	})

	t.Run("Mode 2: Randomized Boundary Test", func(t *testing.T) {
		// Prepare test data for mode 2.
		prepareMode2(t)

		// Verify test data for mode 2.
		verifyMode2(t)

		// Execute accuracy test for mode 2.
		runMode2(t)
	})

	t.Run("Mode 3: Single Node Endurance Test", func(t *testing.T) {
		// Prepare test data for mode 3.
		prepareMode3(t)

		// Verify test data for mode 3.
		verifyMode3(t)

		// Execute accuracy test for mode 3.
		runMode3(t)
	})
}

// shuffleSlice randomly shuffles the elements in the slice.
func shuffleSlice(slice []int64, rng *rand.Rand) {
	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := rng.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}
