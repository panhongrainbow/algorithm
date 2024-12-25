package bptestplan

type NewBpTestDraft struct {
	Goal        string
	Description string

	RandomMaxKey     int64
	RandomMinKey     int64
	RandomTotalCount int64
}

func (b NewBpTestDraft) CreateRandomData() (testData [(2 ^ 63) - 1]int64) {

	for i := 0; i < (2^63)-1; i++ {
		testData[i] = -1
	}

	return
}
