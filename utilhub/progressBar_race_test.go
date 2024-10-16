/*
 * 🏁 Race Condition Testing
 * 📅 Created: 2024-10-16
 * 👤 Author: PanHong
 *
 * 🔍 Test Type: Race Condition Testing
 * 🎯 Goal: Detect race conditions in concurrent access to shared data
 * 📜 Description: Uses the -race option to automatically detect potential data races or conflicts in concurrent code.
 */

package utilhub

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

// ...
func Test_Race_ProgressBar(t *testing.T) {

	// Subtest 1: Collect progress bar messages using ListPrint.
	t.Run("ListPrint", func(t *testing.T) {
		// Initialize a slice to hold the collected progress bar messages.
		var collected []testBarMessage

		// Set total steps and progress bar length.
		var totalSteps uint32 = 10000
		barLength := 50

		// Create a new progress bar.
		progressBar, err := NewProgressBar("Test Progress", totalSteps, barLength)
		assert.NoError(t, err)

		// Create a mutex to protect access to collected messages.
		var mu sync.Mutex

		// Run the ListPrint function concurrently to collect progress bar messages.
		go func() {
			mu.Lock()
			collected = progressBar.ListPrint(T{})
			mu.Unlock()
		}()

		// Use WaitGroup to wait for all goroutines to complete.
		var wg sync.WaitGroup

		// Start multiple goroutines to update the progress bar.
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < int(totalSteps/1000); j++ {
					progressBar.UpdateBar()
					// time.Sleep(10 * time.Millisecond) // Simulate some work.
				}
			}()
		}

		// Wait for all goroutines to complete.
		wg.Wait()

		// Mark the progress bar as complete.
		progressBar.Complete()

		// Wait for the progress bar's printer to stop.
		<-progressBar.WaitForPrinterStop()

		// Collect the collected messages safely.
		mu.Lock()
		// Use require to check if the collected messages count is greater than 0.
		require.Greater(t, len(collected), 0, "Expected at least one message to be collected.")

		mu.Unlock()
	})

	// Subtest 2: Collect progress bar messages using ListenPrinter.
	t.Run("ListenPrinter", func(t *testing.T) {

		// Set total steps and progress bar length.
		var totalSteps uint32 = 10000000
		barLength := 50

		// Create a new progress bar.
		progressBar, err := NewProgressBar("Test Progress", totalSteps, barLength)
		assert.NoError(t, err)

		// Run the ListenPrinter function concurrently to collect progress bar messages.
		go progressBar.ListenPrinter()

		// Use WaitGroup to wait for all goroutines to complete.
		var wg sync.WaitGroup

		// Start multiple goroutines to update the progress bar.
		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < int(totalSteps)/10000; j++ {
					progressBar.UpdateBar()
					// time.Sleep(500 * time.Millisecond) // Simulate some work.
				}
			}()
		}

		// Wait for all goroutines to complete.
		wg.Wait()

		// Mark the progress bar as complete.
		progressBar.Complete()

		// Wait for the progress bar's printer to stop.
		<-progressBar.WaitForPrinterStop()
	})
}
