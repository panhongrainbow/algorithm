package utilhub

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// =====================================================
//                  🛠️ Progress Bar (Tool)
// =====================================================

// ProgressBar ⛏️ struct for managing and tracking progress.
type ProgressBar struct {
	// Basic properties
	name      string // Name of the progress bar.
	total     int    // Total number of steps or units to track.
	barLength int    // The visual length of the progress bar.

	// Tracking progress
	precision        int // Number of decimal places for displaying the progress percentage.
	currentProcess   int // Current progress value.
	lastFilledLength int // Tracks the last filled position to avoid redundant updates.

	// Timezone configuration
	timezone string         // Timezone for displaying the start time of the progress bar.
	location *time.Location // Time.Location object for the specified timezone.

	// Timing information
	startTime time.Time // Start time of the progress tracking.
	endTime   time.Time // End time, set when progress completes.
	complete  bool      // Indicates whether the progress has completed.

	// Time control and synchronization
	updateInterval int // Time interval between each update (in milliseconds).
	// ticker         *time.Ticker // Controls the frequency of updates (regular refreshes).
	ticker <-chan time.Time // Controls the frequency of updates (regular refreshes).

	// Display properties
	barColor   string // ANSI color code for the progress bar display.
	resetColor string // ANSI reset code to revert colors after rendering the progress bar.
}

// BarOption ⛏️ defines a function type for configuring the ProgressBar.
type BarOption func(*ProgressBar)

// WithTracking sets the precision for progress percentage display.
func WithTracking(precision int) BarOption {
	return func(pb *ProgressBar) {
		pb.precision = precision
	}
}

// WithTimeZone sets the timezone string to record the start and end times in the specified timezone.
func WithTimeZone(timeZoneName string) BarOption {
	return func(pb *ProgressBar) {
		pb.timezone = timeZoneName
	}
}

// WithTimeControl sets the ticker for controlling the progress bar update frequency, in milliseconds.
func WithTimeControl(updateInterval int) BarOption {
	return func(pb *ProgressBar) {
		pb.updateInterval = updateInterval
		// pb.ticker.C = time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
		// pb.ticker = time.After(time.Duration(updateInterval) * time.Millisecond)
	}
}

// WithDisplay sets the color for rendering the progress bar using an ANSI color code.
func WithDisplay(color string) BarOption {
	return func(pb *ProgressBar) {
		pb.barColor = color
	}
}

// NewProgressBar ⛏️ initializes and returns a ProgressBar with optional configurations.
func NewProgressBar(name string, total, barLength int, opts ...BarOption) (*ProgressBar, error) {
	// Create a default ProgressBar with the required parameters.
	pb := &ProgressBar{
		// Basic properties
		name:      name,      // Name of the progress bar.
		total:     total,     // Total units to be tracked.
		barLength: barLength, // Length of the progress bar in characters.

		// Tracking progress
		precision:        2, // Default decimal precision for progress percentage.
		currentProcess:   0, // Current progress.
		lastFilledLength: 0, // Track the length of the last update to reduce redundant updates.

		// Timezone configuration
		timezone: "UTC", // Default timezone is set to UTC.
		// location: will be updated (1) (loaded based on the timezone)

		// Timing information
		// startTime: will be updated (2)
		// endTime:   will be updated (3)
		complete: false, // Indicates if the progress bar has completed.

		// Time control and synchronization
		updateInterval: 1000, // Default update interval in milliseconds.
		// ticker: will be updated (4)

		// Display properties
		barColor:   BrightCyan, // Default color for the progress bar.
		resetColor: Reset,      // Reset color to avoid affecting subsequent terminal output.
	}

	// Apply any optional configurations to the default ProgressBar.
	for _, opt := range opts {
		opt(pb)
	}

	// Set the start/end time using the specified timezone.
	loc, err := time.LoadLocation(pb.timezone)
	if err != nil {
		return nil, err
	}
	pb.location = loc // Updated location based on the timezone (1)

	// Set the start time using the specified timezone.
	pb.startTime = time.Now().In(loc) // Start time is set after loading the location (2)

	// If an update interval is provided, initialize the ticker for updates.
	if pb.updateInterval > 0 {
		pb.ticker = time.After(time.Duration(pb.updateInterval) * time.Millisecond)
		// pb.ticker = time.NewTicker(time.Duration(pb.updateInterval) * time.Millisecond) // Initialize the ticker (4)
	}

	return pb, nil
}

// UpdateBar ⛏️ updates the progress bar based on the current count.
func (pb *ProgressBar) UpdateBar() {
	// Return if the progress is already complete.
	if pb.currentProcess == pb.total {
		return
	}

	// Adjust if current process exceeds the total value.
	if pb.currentProcess > pb.total {
		pb.currentProcess = pb.total
		return
	}

	// Increment the current process by one step.
	pb.currentProcess++

	// Calculate the current progress percentage.
	progress := float64(pb.currentProcess) / float64(pb.total)
	filledLength := int(progress * float64(pb.barLength))

	// Format the progress percentage with the specified precision.
	percentage := progress * 100
	if percentage > 100 {
		percentage = 100
	}

	// Format the percentage string using the specified precision.
	format := fmt.Sprintf("%%.%df", pb.precision)
	percentageStr := fmt.Sprintf(format, percentage)

	// Update the progress bar if the filled length has changed.
	if filledLength != pb.lastFilledLength {
	LOOP:
		select {
		case <-pb.ticker:
			// Use "█" to represent the completed portion, and "░" for the remaining portion.
			bar := ""
			for i := 0; i < filledLength; i++ {
				bar += "█"
			}
			for i := filledLength; i < pb.barLength; i++ {
				bar += "░"
			}

			// Print the progress bar with color, along with the percentage.
			if pb.name != "" {
				fmt.Printf("\r%s: %s[%s] %s%%%s", pb.name, pb.barColor, bar, percentageStr, pb.resetColor)
			} else {
				fmt.Printf("\rProgress: %s[%s] %s%%%s", pb.barColor, bar, percentageStr, pb.resetColor)
			}

			// Update the last filled length to avoid redundant updates.
			pb.lastFilledLength = filledLength

			// Reset the ticker for the next update interval.
			pb.ticker = time.After(time.Duration(pb.updateInterval) * time.Millisecond)
		default:
			// Exit the loop if no ticker event occurs.
			break LOOP
		}
	}

	// If progress is complete, stop the ticker.
	if pb.currentProcess == pb.total {
		// pb.ticker.Stop() is not necessary when using time.After.
		pb.ticker = nil
	}
}

// Complete ⛏️ marks the progress bar as complete.
func (pb *ProgressBar) Complete() {
	// Check if the progress bar is already complete.
	if pb.complete == false {
		// Ensure the current progress is less than or equal to the total.
		if pb.currentProcess <= pb.total {
			// Set the end time to the current time in the specified location.
			pb.endTime = time.Now().In(pb.location)

			// Set the current process to the total to mark it as fully completed.
			pb.currentProcess = pb.total

			// Format the progress percentage to display 100%.
			format := fmt.Sprintf("%%.%df", pb.precision)
			percentageStr := fmt.Sprintf(format, 100.0)

			// Generate the progress bar string with all filled blocks.
			bar := ""
			for i := 0; i < int(float64(pb.barLength)); i++ {
				bar += "█"
			}

			// Print the completed progress bar with color and 100% percentage.
			if pb.name != "" {
				fmt.Printf("\r%s: %s[%s] %s%%%s", pb.name, pb.barColor, bar, percentageStr, pb.resetColor)
			} else {
				fmt.Printf("\rProgress: %s[%s] %s%%%s", pb.barColor, bar, percentageStr, pb.resetColor)
			}

			// Mark the progress bar as complete.
			pb.complete = true
		}

		// Set the ticker to nil as no further updates are required.
		pb.ticker = nil

		// Print a newline to indicate completion of the progress bar.
		fmt.Println() // Print a newline to signify that the progress bar is complete.
	}
}

// Report ⛏️ generates and prints a detailed progress report in a formatted table.
func (pb *ProgressBar) Report() error {
	// If the progress is not finished, return an error message.
	if !pb.complete {
		return errors.New("progress is not yet complete")
	}

	// Calculate the total time that has elapsed between the start and the end.
	elapsed := pb.endTime.Sub(pb.startTime)

	// Define fixed widths for the table's fields and values to ensure proper alignment.
	fieldWidth := 20
	valueWidth := 35
	totalWidth := fieldWidth + valueWidth + 7
	// Create a border for the table using a repeated pattern for visual clarity.
	border := BrightYellow + strings.Repeat("=", totalWidth) + Reset
	divider := BrightYellow + strings.Repeat("-", totalWidth) + Reset

	// Print the report title centered within the table, using padding to adjust its position.
	title := "Progress Bar Report"
	titleWidth := len(title)
	padding := (totalWidth - titleWidth) / 2
	fmt.Println(BrightMagenta + border + Reset)
	fmt.Printf("%s|%s%s%s|%s\n", BrightMagenta, strings.Repeat(" ", padding), title, strings.Repeat(" ", padding-1), Reset)
	fmt.Println(BrightMagenta + border + Reset)

	// Print the table header, highlighting the column titles for "Field" and "Value".
	fmt.Printf("%s| %-*s | %-*s |%s\n", BrightRed, fieldWidth, "Field", valueWidth, "Value", Reset)
	fmt.Println(BrightRed + divider + Reset)

	// Print each row of the table with the task's details, formatted to align fields and values.
	fmt.Printf("%s| %-*s | %-*s |%s\n", DarkYellow, fieldWidth, "Task Name", valueWidth, pb.name, Reset) // %-*s ensures left alignment.
	fmt.Printf("%s| %-*s | %-*s |%s\n", DarkYellow, fieldWidth, "Start Time", valueWidth, pb.startTime.Format(time.RFC1123), Reset)
	fmt.Printf("%s| %-*s | %-*s |%s\n", DarkYellow, fieldWidth, "End Time", valueWidth, pb.endTime.Format(time.RFC1123), Reset)
	fmt.Printf("%s| %-*s | %-*s |%s\n", DarkYellow, fieldWidth, "Elapsed Time", valueWidth, elapsed.String(), Reset)
	fmt.Printf("%s| %-*s | %-*d |%s\n", DarkYellow, fieldWidth, "Total Tasks", valueWidth, pb.total, Reset)
	fmt.Printf("%s| %-*s | %-*d |%s\n", DarkYellow, fieldWidth, "Completed Tasks", valueWidth, pb.currentProcess, Reset)

	// Print a closing border to signal the end of the report.
	fmt.Println(BrightMagenta + border + Reset)

	return nil
}
