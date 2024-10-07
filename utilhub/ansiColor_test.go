package utilhub

import (
	"fmt"
	"testing"
)

// Test_AnsiColorOutput prints the colors for manual verification.
func Test_AnsiColorOutput(t *testing.T) {
	colors := []struct {
		name  string // The name of the color for identification.
		color string // The ANSI color code as a string.
	}{
		// Reset color to default.
		{"Reset", Reset},

		// Dark color codes (normal intensity).
		{"DarkBlack", DarkBlack}, // Actually, this represents Black.
		{"DarkRed", DarkRed},
		{"DarkGreen", DarkGreen},
		{"DarkYellow", DarkYellow},
		{"DarkBlue", DarkBlue},
		{"DarkMagenta", DarkMagenta},
		{"DarkCyan", DarkCyan},
		{"DarkWhite", DarkWhite}, // Actually, this represents Gray.

		// Bright color codes (high intensity).
		{"BrightBlack", BrightBlack}, // Actually, this represents Gray.
		{"BrightRed", BrightRed},
		{"BrightGreen", BrightGreen},
		{"BrightYellow", BrightYellow},
		{"BrightBlue", BrightBlue},
		{"BrightMagenta", BrightMagenta},
		{"BrightCyan", BrightCyan},
		{"BrightWhite", BrightWhite}, // Actually, this represents Bright White.
	}

	// Print instructions for manual inspection.
	fmt.Println("Manual Inspection Required:")
	fmt.Println("请进行人工检查:")
	fmt.Println("請進行人工檢查:")

	// Print colors in a table format for manual inspection.
	fmt.Println("Color Code Table:")
	fmt.Println("----------------------------")
	fmt.Printf("| %-15s | %-6s |\n", "Color", "Sample")
	fmt.Println("----------------------------")
	for _, c := range colors {
		fmt.Printf("| %-15s | %-15s |\n", c.name, fmt.Sprintf("%sSample%s", c.color, Reset))
	}
	fmt.Println("----------------------------")
}
