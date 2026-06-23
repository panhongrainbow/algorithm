package randomizedBoundary

import (
	"testing"

	"github.com/panhongrainbow/go-algorithm/utilhub"
	"github.com/stretchr/testify/require"
)

// Test_Model2_Generate_Check_RandomSet verifies BpTestRandomizedBoundary's random data generation by setting a count of 50,
// creating a dataset, and validating its integrity through the model's own checks.
func Test_Model2_Generate_Check_RandomSet(t *testing.T) {
	// Set the total count for random data generation to 50 in pool.
	utilhub.SetAutoRandomTotalCount(50)

	// Verify that the random total count was correctly set to 50.
	require.Equal(t, utilhub.GetAutoRandomTotalCount(), int64(50))

	// Create an instance of BpTestRandomizedBoundary.
	bptest2 := &BpTestRandomizedBoundary{}

	// Generate a random test dataset using the model.
	testDataSet, err := bptest2.GenerateRandomSet()
	require.NoError(t, err, "failed to generate test data")

	// Check the validity of the generated random dataset.
	err = bptest2.CheckRandomSet(testDataSet)
	require.NoError(t, err, "failed to check test data")

	// Force reload the configuration to reset any changes made during testing.
	utilhub.ForceReloadAutoTest()
}
