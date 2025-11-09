package bptestModel1

import (
	"testing"

	"github.com/panhongrainbow/algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

// Test_Model1_Generate_Check_RandomSet verifies BpTestModel1's random data generation by setting a count of 50,
// creating a dataset, and validating its integrity through the model's own checks.
func Test_Model1_Generate_Check_RandomSet(t *testing.T) {
	// Set the total count for random data generation to 50 in pool.
	utilhub.SetRandomTotalCount(100)

	// Verify that the random total count was correctly set to 50.
	require.Equal(t, utilhub.GetRandomTotalCount(), int64(100))

	// Create an instance of BpTestModel1.
	bptest1 := &BpTestModel1{}

	// Generate a random test dataset using the model.
	testDataSet, err := bptest1.GenerateRandomSet(1, 30)
	require.NoError(t, err, "failed to generate test data")

	// Check the validity of the generated random dataset.
	err = bptest1.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to check test data")

	// Force reload the configuration to reset any changes made during testing.
	utilhub.ForceReloadConfig()
}
