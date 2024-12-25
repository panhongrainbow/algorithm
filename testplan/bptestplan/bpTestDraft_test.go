package bptestplan

import "testing"

func TestCreateRandomData(t *testing.T) {
	b := NewBpTestDraft{}

	testData := b.CreateRandomData()

	// Check if the length of testData is correct
	if len(testData) != (2^63)-1 {
		t.Errorf("Expected length of testData to be %d, but got %d", (2^63)-1, len(testData))
	}

	// Check if all elements in testData are -1
	for i := range testData {
		if testData[i] != -1 {
			t.Errorf("Expected testData[%d] to be -1, but got %d", i, testData[i])
		}
	}
}