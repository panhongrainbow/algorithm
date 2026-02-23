package utilhub

// TestTimeString gets the current time as a formatted string in the given time zone.
func TestTimeString(format string, timeZone string) string {
	// Call the function GetNowTimeString from the utilhub package to get the current time in string format.
	str, err := GetNowTimeString(format, timeZone)
	if err != nil {
		// If an error occurs, panic and terminate the program.
		panic(err)
	}
	// Return the formatted time string.
	return str
}
