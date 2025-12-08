package model3

import (
	bptestModel "github.com/panhongrainbow/algorithm/testdata/share"
)

type BpTestModel3 struct{}

func (model3 *BpTestModel3) GenerateRandomSet() ([]int64, error) {
	model := bptestModel.BpTestShare{}
	return model.ShareGenerateRandomSet(5)
}

// CheckRandomSet 🧮 checks the validity of a random data set by comparing the positive and negative numbers.
func (model3 *BpTestModel3) CheckRandomSet(dataSet []int64) error {
	model := bptestModel.BpTestShare{}
	return model.CheckRandomSet(dataSet)
}
