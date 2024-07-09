package bpTree

import (
	"errors"
	"fmt"
	"sort"
)

// delAndDir performs data deletion based on automatic direction detection.  // 这是 B 加树的方向性删除入口
// 自动判断资料删除方向，其實會由不同方向進行刪除

/*
 为何要先优先向左删除资料，因最左边的相同值被删除时，就会被后面相同时递补，比较不会更动到边界值 ✌️
*/

func (inode *BpIndex) delAndDir(item BpItem) (deleted, updated bool, ix int, edgeValue int64, err error) {
	// 搜寻 🔍 (最右边 ➡️)
	// Use binary search to find the index (ix) where the key should be deleted.
	ix = sort.Search(len(inode.Index), func(i int) bool {
		return inode.Index[i] > item.Key // 一定要大于，所以会找到最右边 ‼️
	})

	// 决定 ↩️ 是否要向左
	// Check if deletion should be performed by the leftmost node first.
	if len(inode.Index) > 0 && len(inode.IndexNodes) > 0 &&
		(ix-1) >= 1 && len(inode.IndexNodes)-1 >= (ix-1) { // 如果当前节点的左边有邻居

		// If it is continuous data (same value) (5❌ - 5 - 5 - 5 - 5 - 6 - 7 - 8)
		length := len(inode.IndexNodes[ix-1].Index) // 为了左边邻居节点最后一个索引值
		if len(inode.IndexNodes) > 0 &&             // 预防 panic 的检查
			len(inode.IndexNodes[ix].Index) > 0 && len(inode.IndexNodes[ix-1].Index) > 0 && // 预防 panic 的检查
			length > 0 && inode.IndexNodes[ix].Index[0] == inode.IndexNodes[ix-1].Index[length-1] { // 最后决定，如果最接近的索引节点有相同的索引值 ‼️

			// 搜寻 🔍 (最左边 ⬅️) (一切重来，重头开始向左搜寻)
			deleted, updated, ix, err = inode.deleteToLeft(item) // Delete to the leftmost node ‼️ (向左砍)

			// 中断了，不再考虑向右搜寻 ⚠️
			return
		}
	}

	// 搜寻 🔍 (最右边 ➡️)
	// If it is discontinuous data (different values) (5 - 5 - 5 - 5 - 5❌ - 6 - 7 - 8)
	deleted, updated, edgeValue, _, ix, err = inode.deleteToRight(item) // Delete to the rightmost node ‼️ (向右砍)

	// Return the results.
	return
}

// deleteToRight is designed to delete from the rightmost side within continuous data.  (5 - 5 - 5 - 5 - 5❌ - 6 - 7 - 8)

// deleteToRight 先放前面，因为 deleteToLeft 会抄 deleteToRight 的内容
func (inode *BpIndex) deleteToRight(item BpItem) (deleted, updated bool, edgeValue int64, status int, ix int, err error) {
	// Initialize the return value first.
	status = edgeValueInit
	edgeValue = -1

	// ✈️ Process the Index Node.
	if len(inode.IndexNodes) > 0 {
		ix = sort.Search(len(inode.Index), func(i int) bool {
			// 🖍️ The `Sort` function stops when the condition is met.
			// When it equals, it meets the condition later, so it will delete the data on the far right.
			// When it is greater than or equal to, it meets the condition earlier, so it will delete the data on the far left.
			return inode.Index[i] > item.Key // 在最右边 ‼️
		})

		// Entering the Recursive Function. 🔁
		deleted, updated, edgeValue, status, _, err = inode.IndexNodes[ix].deleteToRight(item)

		// Mechanism for updating edge values.
		if ix > 0 && status == edgeValueUpload {
			// 🖍️ In this block, the edge values will be uploaded.
			// When uploaded to a location where ix is greater than 0, it becomes an index and stops uploading.
			// (边界值会变成索引并中止)

			inode.Index[ix-1] = edgeValue
			updated = false
			status = edgeValueInit

			// Interrupted, index updated, no uploading. ⚠️
			return
		} else if ix == 0 && status == edgeValueUpload {
			// 🖍️ When uploaded to a location where ix equals 0, it continues to upload immediately until the boundary value is not 0.
			// (IX 为 0 时不停上传)

			// Continuous uploading. ⚠️
			return
		}

		// 🖍️ In this block, (temporarily) decide whether you want to update the boundary values or upload the

		// 🖐️ 状态变化 [LeaveBottom] -> Any
		if status == edgeValueRenew {

			// ⚠️ 状况一 用边界值去更新任意索引

			// 🖐️ 状态变化 [LeaveBottom] -> [Init]
			// 看到 LeaveBottom 状态时，就代表准备要更新边界值，但更新的索引不一定在最左边
			if ix-1 >= 0 {

				fmt.Println(">>>>> 更新完成")

				inode.Index[ix-1] = edgeValue

				status = edgeValueInit
				return
			} else {
				fmt.Println(">>>>> 进行上传")
				status = edgeValueUpload
				return
			}

			// To make temporary corrections, mainly to identify the problems.
		} else { // if status == statusBorrowFromIndexNode || item.Key == 960 {

			if status != statusBorrowFromIndexNode && inode.IndexNodes[ix].DataNodes != nil {
				_, _, edgeValue, err, status = inode.borrowFromBottomIndexNode(ix)
				return
			}

			if status != statusBorrowFromIndexNode && len(inode.IndexNodes[ix].Index) == 0 {
				if len(inode.IndexNodes[ix].Index) == 0 {
					inode.IndexNodes[ix].Index = []int64{edgeValue}
				}

				ix, edgeValue, status, err = inode.borrowFromIndexNode(ix)
				if ix == 0 && status == edgeValueChanges {
					status = edgeValueUpload
					return
				}
				return
			}

			if status == statusBorrowFromIndexNode {
				if len(inode.IndexNodes[ix].Index) == 0 {
					inode.IndexNodes[ix].Index = []int64{edgeValue}
				}

				ix, edgeValue, status, err = inode.borrowFromIndexNode(ix)
				if ix == 0 && status == edgeValueChanges {
					status = edgeValueUpload
					return
				}
			}

			return
		}

		// If the index at position ix becomes invalid. ‼️
		// 删除导致锁引失效 ‼️
		if len(inode.IndexNodes[ix].Index) == 0 { // invalid ❌
			if len(inode.IndexNodes[ix].DataNodes) >= 2 { // DataNode 🗂️

				// 之后从这开始开发 ‼️

				var borrowed bool

				borrowed, _, edgeValue, err, status = inode.borrowFromBottomIndexNode(ix) // Will borrow part of the node (借结点). ‼️  // 🖐️ for index node 针对索引节点
				// 看看有没有向索引节点借到资料

				if err != nil && !errors.Is(err, fmt.Errorf("the index is still there; there is no need to borrow nodes")) {
					return
				}

				if borrowed == true { // 当向其他索引节点借完后，在执行 borrowFromIndexNode，重新计算边界值

					if ix == 0 && status == edgeValueChanges {
						fmt.Println(">>>>> 进行上传")
						status = edgeValueUpload
						return
					}

					if len(inode.IndexNodes) > 0 && // 预防性检查
						len(inode.IndexNodes[0].DataNodes) > 0 && // 预防性检查
						len(inode.IndexNodes[0].DataNodes[0].Items) > 0 { // 预防性检查

						edgeValue = inode.IndexNodes[0].DataNodes[0].Items[0].Key // 边界值是由 索引节点中取出，所以可以直接把边界值放入 索引  ‼️‼️

						if edgeValue != -1 && len(inode.Index) == 0 { // 如果有正确取得 边界值 后
							inode.Index = []int64{edgeValue}
							status = statusBorrowFromIndexNode
						}
					}

					return
				}
			}
		}

		// Return the results of the deletion.
		return
	}

	// ✈️ Process the Data Node.
	if len(inode.DataNodes) > 0 {
		// Call the deleteBottomItem method on the current node as it is close to the bottom layer.
		// This signifies the beginning of deleting data. (接近资料层)

		// Here, adjustments may be made to IX (IX 在这里可能会被修改) ‼️
		// var edgeValue int64
		deleted, updated, ix, edgeValue, status = inode.deleteBottomItem(item) // 🖐️ for data node 针对资料节点
		if ix == 0 && status == edgeValueChangesByDelete {                     // 当 ix 为 0 时，才要处理边界值的问题 (ix == 0，是特别加入的)
			status = edgeValueRenew
		}

		// The individual data node is now empty, and
		// it is necessary to start borrowing data from neighboring nodes.
		if len(inode.DataNodes[ix].Items) == 0 { // 会有一边的资料节点没有任何资料
			var borrowed bool
			borrowed, edgeValue, err = inode.borrowFromDataNode(ix) // Will borrow part of the data node. (向资料节点借资料)

			// 先检查是否有错误
			if err != nil {
				status = statusError
				return
			}

			if edgeValue != -1 {
				status = edgeValueUpload
				return
			}

			// If the data node cannot be borrowed, then information should be borrowed from the index node later.
			// 如果资料节点借到，就不需后续处理
			if borrowed == true {
				updated = true
				return
			}

			// 如果使用 borrowFromDataNode 没有借到资料，就要进行以下处理 ‼️ ‼️

			// ⚠️ 状况一 索引节点资料过少，整个节点失效
			// During the deletion process, the node's index may become invalid.
			// 如果资料节点数量过少
			if len(inode.DataNodes) <= 2 { // 资料节点数量过少

				inode.Index = []int64{}

				// 状况更新
				updated = true

				// 有可能新增状态
				if item.Key == 960 {
					fmt.Println()
				}
				fmt.Println("有可能新增状态")

				// 直接中断
				return
			}

			// ⚠️ 状况二 索引节点有一定数量的资料，删除部份资料后，还能维持为一个节点
			// Wipe out the empty data node at the specified 'ix' position directly.
			// 如果资料节点删除资料后，还是维持为一个节点的定义，就要进行抹除部份 ix 位置上的资料 ‼️
			if len(inode.Index) != 0 {
				// Rebuild the connections between data nodes.
				if inode.DataNodes[ix].Previous == nil {
					inode.DataNodes[ix].Next.Previous = nil

					status = edgeValueInit
				} else if inode.DataNodes[ix].Next == nil {
					inode.DataNodes[ix].Previous.Next = nil

					status = edgeValueInit
				} else {
					inode.DataNodes[ix].Previous.Next = inode.DataNodes[ix].Next
					inode.DataNodes[ix].Next.Previous = inode.DataNodes[ix].Previous

					status = edgeValueInit
				}

				// Reorganize nodes.
				if ix != 0 {
					inode.Index = append(inode.Index[:ix-1], inode.Index[ix:]...)             // Erase the position of ix - 1.
					inode.DataNodes = append(inode.DataNodes[:ix], inode.DataNodes[ix+1:]...) // Erase the position of ix.

					status = edgeValueInit
				} else if ix == 0 { // Conditions have already been established earlier, with the index length not equal to 0. ‼️
					inode.Index = inode.Index[1:]
					inode.DataNodes = inode.DataNodes[1:]

					// 边界值要立刻进行修改
					edgeValue = inode.DataNodes[0].Items[0].Key
					status = edgeValueUpload
				}
			}
		}

	}

	// Return the results of the deletion.
	return
}

// deleteToLeft is a method of the BpIndex type that deletes the leftmost specified BpItem. (由左边删除 👈 ‼️)
func (inode *BpIndex) deleteToLeft(item BpItem) (deleted, updated bool, ix int, err error) {
	panic("Currently developing right deletion, not developing left deletion.")

	// Return the results of the deletion.
	return
}
