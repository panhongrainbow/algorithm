package bpTree

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func loadBasicDeletionExample() (basicDeletionBpTree *BpTree) {
	// Generate continuous data for insertion.
	var basicDeletionNumbers = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}

	// Initialize B Plus tree.
	basicDeletionBpTree = NewBpTree(4)

	// Start inserting data.
	for i := 0; i < len(basicDeletionNumbers); i++ {
		// Insert data entries continuously.
		basicDeletionBpTree.InsertValue(BpItem{Key: basicDeletionNumbers[i]})
	}

	// Complete this function.
	return
}

// Test_Check_borrowFromDataNode_Function is primarily used to test the borrowFromDataNode function.
// More details can be found in Chapter 2.3.1 `Borrow from Neighbor` in the documentation.
func Test_Check_borrowFromDataNode_Function(t *testing.T) {
	// 🧪 This test is mainly used to test the scenario of Status 1.
	t.Run("Status 1 in Chapter 2.3.1", func(t *testing.T) {
		// Load a simple B Plus Tree where max degree is 4.
		basicDeletionBpTree := loadBasicDeletionExample()

		// Deleting the Non-Edge-Value 14.
		deleted, _, _, _ := basicDeletionBpTree.RemoveValue(BpItem{Key: 14})
		require.True(t, deleted)

		// Deleting the Non-Edge-Value 13.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 13})
		require.True(t, deleted)

		// Check the index node of the first level after deleting data.
		require.Equal(t, []int64{7, 15}, basicDeletionBpTree.root.Index)

		// Check the index node of the second level after deleting data.
		require.Equal(t, []int64{16, 17, 19}, basicDeletionBpTree.root.IndexNodes[2].Index)

		// Check the data nodes of the third level after deleting data.
		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[0].Items))
		require.Equal(t, int64(15), basicDeletionBpTree.root.IndexNodes[2].DataNodes[0].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[1].Items))
		require.Equal(t, int64(16), basicDeletionBpTree.root.IndexNodes[2].DataNodes[1].Items[0].Key)
	})

	// 🧪 This test is mainly used to test the scenario of Status 2-1.
	t.Run("Status 2-1 in Chapter 2.3.1", func(t *testing.T) {
		// Load a simple B Plus Tree where max degree is 4.
		basicDeletionBpTree := loadBasicDeletionExample()

		// Deleting the Non-Edge-Value 20.
		deleted, _, _, _ := basicDeletionBpTree.RemoveValue(BpItem{Key: 20})
		require.True(t, deleted)

		// Deleting the Inner-Edge-Value 19.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 19})
		require.True(t, deleted)

		// Deleting the inner-Edge-Value 21.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 21})
		require.True(t, deleted)

		// Check the index node of the first level after deleting data.
		require.Equal(t, []int64{7, 13}, basicDeletionBpTree.root.Index)

		// Check the index node of the second level after deleting data.
		require.Equal(t, []int64{15, 17, 18}, basicDeletionBpTree.root.IndexNodes[2].Index)

		// Check the data nodes of the third level after deleting data.
		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items))
		require.Equal(t, int64(17), basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items))
		require.Equal(t, int64(18), basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items[0].Key)
	})

	// 🧪 This test is mainly used to test the scenario of Status 2-2.
	t.Run("Status 2-2 in Chapter 2.3.1", func(t *testing.T) {
		// Load a simple B Plus Tree where max degree is 4.
		basicDeletionBpTree := loadBasicDeletionExample()

		// Deleting the Non-Edge-Value 20.
		deleted, _, _, _ := basicDeletionBpTree.RemoveValue(BpItem{Key: 20})
		require.True(t, deleted)

		// Deleting the Inner-Edge-Value 19.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 19})
		require.True(t, deleted)

		// Deleting the inner-Edge-Value 18.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 18})
		require.True(t, deleted)

		// Deleting the inner-Edge-Value 17.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 17})
		require.True(t, deleted)

		// Check the index node of the first level after deleting data.
		require.Equal(t, []int64{7, 13}, basicDeletionBpTree.root.Index)

		// Check the index node of the second level after deleting data.
		require.Equal(t, []int64{15, 16, 21}, basicDeletionBpTree.root.IndexNodes[2].Index)

		// Check the data nodes of the third level after deleting data.
		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[1].Items))
		require.Equal(t, int64(15), basicDeletionBpTree.root.IndexNodes[2].DataNodes[1].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items))
		require.Equal(t, int64(16), basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items))
		require.Equal(t, int64(21), basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items[0].Key)
	})

	// 🧪  This test is mainly used to test the scenario of Status 3.
	t.Run("Status 3 in Chapter 2.3.1", func(t *testing.T) {
		// Load a simple B Plus Tree where max degree is 4.
		basicDeletionBpTree := loadBasicDeletionExample()

		// Deleting the Non-Edge-Value 18.
		deleted, _, _, _ := basicDeletionBpTree.RemoveValue(BpItem{Key: 18})
		require.True(t, deleted)

		// Deleting the Inner-Edge-Value 17.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 17})
		require.True(t, deleted)

		// Check the index node of the first level after deleting data.
		require.Equal(t, []int64{7, 13}, basicDeletionBpTree.root.Index)

		// Check the index node of the second level after deleting data.
		require.Equal(t, []int64{15, 19, 20}, basicDeletionBpTree.root.IndexNodes[2].Index)

		// Check the data nodes of the third level after deleting data.
		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items))
		require.Equal(t, int64(19), basicDeletionBpTree.root.IndexNodes[2].DataNodes[2].Items[0].Key)

		require.Equal(t, 2, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items))
		require.Equal(t, int64(20), basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items[0].Key)
		require.Equal(t, int64(21), basicDeletionBpTree.root.IndexNodes[2].DataNodes[3].Items[1].Key)
	})
}

// Test_Check_borrowFromBottomIndexNode_Function will verify the following process:
// it will borrow data from lower-level index nodes.
// However, the process is complex, with at least six scenarios that need to be analyzed one by one.
func Test_Check_borrowFromBottomIndexNode_Function(t *testing.T) {
	// 🧪  This test is mainly used to test the scenario of Status 3.
	t.Run("Scenario in Chapter 2.3.2", func(t *testing.T) {
		// Load a simple B Plus Tree where max degree is 4.
		basicDeletionBpTree := loadBasicDeletionExample()

		// Deleting the Inner-Edge-Value 7.
		deleted, _, _, _ := basicDeletionBpTree.RemoveValue(BpItem{Key: 7})
		require.True(t, deleted)

		// Deleting the new Inner-Edge-Value 8.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 8})
		require.True(t, deleted)

		// Deleting the new Inner-Edge-Value 9.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 9})
		require.True(t, deleted)

		// Deleting the new Inner-Edge-Value 10.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 10})
		require.True(t, deleted)

		// Deleting the new Inner-Edge-Value 11.
		deleted, _, _, _ = basicDeletionBpTree.RemoveValue(BpItem{Key: 11})
		require.True(t, deleted)

		// Check the index node of the first level after deleting data.
		require.Equal(t, []int64{12, 14}, basicDeletionBpTree.root.Index)

		// Check the index node of the second level after deleting data.
		require.Equal(t, []int64{13}, basicDeletionBpTree.root.IndexNodes[1].Index)
		require.Equal(t, []int64{15, 17, 19}, basicDeletionBpTree.root.IndexNodes[2].Index)

		// Check the data nodes of the third level after deleting data.
		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[1].DataNodes[0].Items))
		require.Equal(t, int64(12), basicDeletionBpTree.root.IndexNodes[1].DataNodes[0].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[1].DataNodes[1].Items))
		require.Equal(t, int64(13), basicDeletionBpTree.root.IndexNodes[1].DataNodes[1].Items[0].Key)

		require.Equal(t, 1, len(basicDeletionBpTree.root.IndexNodes[2].DataNodes[0].Items))
		require.Equal(t, int64(14), basicDeletionBpTree.root.IndexNodes[2].DataNodes[0].Items[0].Key)
	})
}
