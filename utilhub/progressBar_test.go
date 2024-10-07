package utilhub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// testBarMessage holds the details of each progress bar message.
type testBarMessage struct {
	humanEyeSpeedInterval float64 // The interval between messages, calculated as a multiple of 10 milliseconds, relative to human eye speed.
	filledLength          int     // The number of filled units in the progress bar.
	percentage            float64 // The current progress percentage (0 to 100%).
}

// ListPrint captures messages printed by the progress bar and stores them for testing purposes.
// It also calculates the time interval between messages.
func (pb *ProgressBar) ListPrint(tt T) (collected []testBarMessage) {
	// Ensure that this function can only be called in a test environment.
	tt.EnsureTestEnvironment()

	// Initialize a variable to track the last message's timestamp.
	var lastTime time.Time

	// Iterate over messages received through the progress bar's print channel.
	for msg := range pb.printChannel {
		// Create a new testBarMessage to store the current message details.
		testMsg := testBarMessage{
			filledLength: msg.filledLength, // Store the filled length.
			percentage:   msg.percentage,   // Store the percentage completion.
		}

		// Get the current timestamp.
		currentTime := time.Now()

		// If lastTime is set (i.e., not zero), calculate the time interval since the last message.
		if !lastTime.IsZero() {
			interval := currentTime.Sub(lastTime)
			intervalMs := interval.Seconds() * 1000         // Convert interval to milliseconds.
			testMsg.humanEyeSpeedInterval = intervalMs / 10 // Convert to a multiple of 10 milliseconds.
		}

		// Update lastTime to the current time for the next iteration.
		lastTime = currentTime

		// Collect the message into the list of messages.
		collected = append(collected, testMsg)
	}

	// Notify that the progress bar has finished processing.
	pb.finishBar <- struct{}{}

	// Return the collected list of messages.
	return
}

// validateCollectedMessages checks if the collected messages meet specific conditions.
// It returns a list of errors for any failed conditions.
func validateCollectedMessages(tt T, collected []testBarMessage, humanEyeSpeedMultiplier float64) []error {
	// Ensure this function is only executed in test mode.
	tt.EnsureTestEnvironment()

	// Create a list to hold any validation errors.
	var errs []error

	// Check the time intervals between messages, starting from the second message.
	for i := 1; i < len(collected)-1; i++ {
		if collected[i].humanEyeSpeedInterval <= humanEyeSpeedMultiplier {
			// If the time interval is too short, add an error.
			errs = append(errs, fmt.Errorf("time interval check failed: message %d has interval %.2f, should be greater than %.2f", i, collected[i].humanEyeSpeedInterval, humanEyeSpeedMultiplier))
		}
	}

	// Check the progress values (filledLength and percentage).
	for i := 1; i < len(collected); i++ {
		// Ensure that the filledLength increases with each message.
		if collected[i].filledLength < collected[i-1].filledLength {
			errs = append(errs, fmt.Errorf("progress check failed: message %d has filledLength %d, should be greater than the previous message's %d", i, collected[i].filledLength, collected[i-1].filledLength))
		}
		// Ensure that the percentage increases with each message.
		if collected[i].percentage < collected[i-1].percentage {
			errs = append(errs, fmt.Errorf("progress check failed: message %d has percentage %.2f, should be greater than the previous message's %.2f", i, collected[i].percentage, collected[i-1].percentage))
		}
	}

	// Return the list of errors (if any).
	return errs
}

// Test_ProcessBar is a unit test for the progress bar functionality.
// It verifies that the progress bar updates as expected and checks the validity of collected messages.

/*
This table effectively outlines how to set the speed for collecting messages from the progress bar based on different update intervals and working rates.

| Test Scenario                              | Update Interval             | Working Rate             | Update Count      | Expected Collected Messages      | Description**                                   |
|--------------------------------------------|-----------------------------|--------------------------|-------------------|----------------------------------|---------------------------------------------------|
| Update Interval Faster Than Working Rate   | 100 milliseconds            | 150 milliseconds         | 10                | 10                               | The progress bar generates a message for each of the 10 updates. |
| Update Interval Slower Than Working Rate   | 150 * 10 + 500 milliseconds | 150 milliseconds         | 10                | 1                                | Due to a longer update interval, only the final completion message is recorded despite 10 updates. |
*/

// When calling progressBar.Complete(), it automatically moves to a new line to allow for report generation. (自动换行，以便生成报告)

func Test_ProcessBar(t *testing.T) {
	t.Run("The Update Interval is Faster Than the Working Rate.", func(t *testing.T) {
		// Initialize a slice to hold the collected progress bar messages.
		var collected []testBarMessage

		// Create a ProgressBar with various configuration options.
		progressBar, _ := NewProgressBar("Download", 10, 70,
			WithDisplay(BrightCyan),     // Set the progress bar display color to Bright Cyan.
			WithTracking(1),             // Enable tracking with a tracking ID of 1.
			WithTimeControl(100),        // Set the update interval to 100 milliseconds.
			WithTimeZone("Asia/Taipei"), // Set the time zone to Asia/Taipei.
		)

		// Run the ListPrint function concurrently to collect progress bar messages.
		go func() {
			collected = progressBar.ListPrint(T{})
		}()

		// Simulate progress by updating the bar 10 times, with a 150-millisecond pause between updates.
		for i := 0; i < 10; i++ {
			progressBar.UpdateBar()
			time.Sleep(150 * time.Millisecond) // Simulate some work being done.
		}

		// Mark the progress bar as complete.
		progressBar.Complete()

		// Wait for the progress bar's printer to stop.
		<-progressBar.WaitForPrinterStop()

		// Validate the collected messages using the validateCollectedMessages function.
		errs := validateCollectedMessages(T{}, collected, 15) // 15 倍的人类眼速

		// If there are errors, use assert.Fail to report all errors.
		for _, err := range errs {
			assert.Fail(t, err.Error())
		}

		// Use assert to check that the correct number of messages were collected.
		assert.Equal(t, 10, len(collected), "Expected 10 collected messages, but got %d", len(collected))
	})
	t.Run("The Update Interval is Slower Than the Working Rate.", func(t *testing.T) {
		// Initialize a slice to hold the collected progress bar messages.
		var collected []testBarMessage

		// Create a ProgressBar with various configuration options.
		progressBar, _ := NewProgressBar("Download", 10, 70,
			WithDisplay(BrightCyan),     // Set the progress bar display color to Bright Cyan.
			WithTracking(1),             // Enable tracking with a tracking ID of 1.
			WithTimeControl(150*10+500), // Set the update interval to 100 milliseconds.
			WithTimeZone("Asia/Taipei"), // Set the time zone to Asia/Taipei.
		)

		// Run the ListPrint function concurrently to collect progress bar messages.
		go func() {
			collected = progressBar.ListPrint(T{})
		}()

		// Simulate progress by updating the bar 10 times, with a 150-millisecond pause between updates.
		for i := 0; i < 10; i++ {
			progressBar.UpdateBar()
			time.Sleep(150 * time.Millisecond) // Simulate some work being done.
		}

		// Mark the progress bar as complete.
		progressBar.Complete()

		// Wait for the progress bar's printer to stop.
		<-progressBar.WaitForPrinterStop()

		// Validate the collected messages using the validateCollectedMessages function.
		errs := validateCollectedMessages(T{}, collected, 200) // 200 倍的人类眼速

		// If there are errors, use assert.Fail to report all errors.
		for _, err := range errs {
			assert.Fail(t, err.Error())
		}

		// Use assert to check that the correct number of messages were collected.
		assert.Equal(t, 1, len(collected), "Expected 10 collected messages, but got %d", len(collected))
	})

}
