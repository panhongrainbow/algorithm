package model3

import (
	bptestModel2 "github.com/panhongrainbow/algorithm/testdata/model2"
)

type BpTestModel3 struct{}

func (model3 *BpTestModel3) GenerateRandomSet() ([]int64, error) {
	model2 := bptestModel2.BpTestModel2{}
	return model2.ShareGenerateRandomSet(5)
}
