package singleNodeEndurance

import (
	bptestModel "github.com/panhongrainbow/go-algorithm/testdata/share"
)

type BpTestSingleNodeEndurance struct{}

func (model3 *BpTestSingleNodeEndurance) GenerateRandomSet() ([]int64, error) {
	model := bptestModel.BpTestShare{}
	return model.ShareGenerateRandomSet(5)
}

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers.
func (model3 *BpTestSingleNodeEndurance) CheckRandomSet(dataSet []int64) error {
	model := bptestModel.BpTestShare{}
	return model.CheckRandomSet(dataSet)
}
