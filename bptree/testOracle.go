package bpTree

import "fmt"

func (tree *BpTree) CheckAndSwapRightContinuity2() error {
	next := tree.root.BpDataHead()
	var temp int64

	for next != nil {
		for i := 0; i < len(next.Items); i++ {
			temp = next.Items[i].Key
			next.Items[i].Key = 0

			right := next.Next
			for right != nil {
				for j := 0; j < len(right.Items); j++ {
					if right.Items[j].Key == 0 {
						tree.root.Print()
						fmt.Println(temp)
						panic("error")
					}
				}
				right = right.Next
			}

			next.Items[i].Key = temp
		}
		next = next.Next
	}
	return nil
}

func (tree *BpTree) CheckAndSwapLeftContinuity2() error {
	previous := tree.root.BpDataTail()
	var temp int64

	for previous != nil {
		for i := 0; i < len(previous.Items); i++ {
			temp = previous.Items[i].Key
			previous.Items[i].Key = 0

			left := previous.Previous
			for left != nil {
				for j := 0; j < len(left.Items); j++ {
					if left.Items[j].Key == 0 {
						tree.root.Print()
						fmt.Println(temp)
						panic("error")
					}
				}
				left = left.Previous
			}

			previous.Items[i].Key = temp
		}
		previous = previous.Previous
	}
	return nil
}

func (tree *BpTree) CheckAndSwapRightContinuity() error {
	next := tree.root.BpDataHead()
	var temp int64
	var count int

	for next != nil {
		for i := 0; i < len(next.Items); i++ {
			temp = next.Items[i].Key
			next.Items[i].Key = 0

			right := next
			for right != nil {
				for j := 0; j < len(right.Items); j++ {
					if right.Items[j].Key == 0 {
						count++
						if count >= 2 {
							tree.root.Print()
							fmt.Println(temp)
							panic("error")
						}
					}
				}
				right = right.Next
			}

			next.Items[i].Key = temp
			count = 0
		}
		next = next.Next
	}
	return nil
}

func (tree *BpTree) CheckAndSwapLeftContinuity() error {
	previous := tree.root.BpDataTail()
	var temp int64
	var count int

	for previous != nil {
		for i := 0; i < len(previous.Items); i++ {
			temp = previous.Items[i].Key
			previous.Items[i].Key = 0

			left := previous
			for left != nil {
				for j := 0; j < len(left.Items); j++ {
					if left.Items[j].Key == 0 {
						count++
						if count >= 2 {
							tree.root.Print()
							fmt.Println(temp)
							panic("error")
						}
					}
				}
				left = left.Previous
			}

			previous.Items[i].Key = temp
			count = 0
		}
		previous = previous.Previous
	}
	return nil
}
