package utilhub

// =====================================================================================================================
//                  🛠️ ANSI Color (Tool)
// ANSI color codes for foreground (text color) and background (background color).
// =====================================================================================================================

// ANSI color codes for foreground (text color).
const (
	// Reset ⛏️ all attributes (color, formatting, etc.)
	Reset = "\033[00m"

	// DarkBlack ... ⛏️ : Dark colors (normal intensity).
	DarkBlack   = "\033[30m" // Black
	DarkRed     = "\033[31m" // Red
	DarkGreen   = "\033[32m" // Green
	DarkYellow  = "\033[33m" // Yellow
	DarkBlue    = "\033[34m" // Blue
	DarkMagenta = "\033[35m" // Magenta
	DarkCyan    = "\033[36m" // Cyan
	DarkWhite   = "\033[37m" // White (actually gray)

	// BrightBlack ... ⛏️ : Bright colors (high intensity).
	BrightBlack   = "\033[90m" // Bright Black (Gray)
	BrightRed     = "\033[91m" // Bright Red
	BrightGreen   = "\033[92m" // Bright Green
	BrightYellow  = "\033[93m" // Bright Yellow
	BrightBlue    = "\033[94m" // Bright Blue
	BrightMagenta = "\033[95m" // Bright Magenta
	BrightCyan    = "\033[96m" // Bright Cyan
	BrightWhite   = "\033[97m" // Bright White
)
