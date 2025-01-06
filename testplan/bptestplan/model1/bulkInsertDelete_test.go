package bptestModel1

import (
	"fmt"
	"testing"
)

func Test_Mode1(t *testing.T) {

	model1 := &BpTestModel1{RandomTotalCount: 10}

	test, err := model1.GenerateRandomSet(1, 11)
	fmt.Println(err)

	fmt.Println(test)

	var sum int64

	for i := 0; i < len(test); i++ {
		sum = sum + test[i]
	}

	fmt.Println(sum)
}
