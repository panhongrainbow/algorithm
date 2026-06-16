package utilhub

import (
	"path/filepath"
)

// ParseToggle ⛏️ loads the toggle configuration from struct tags and applies it to the provided struct.
func ParseToggle(cfg ToggleConfig) error {
	// Prepare the variable outside the closure function.
	var err error
	var projectPath, file string

	// Use Golang's sync.Once to prevent the setting from being overwritten.
	_onesToggleConfig.Do(func() {

		// Get the default configuration directory.
		projectPath, err = GetProjectDir(filepath.Join(ProjectName))
		if err != nil {
			return
		}

		// Get the struct name to use as the filename.
		file, err = GetDefaultStructName(&cfg)
		if err != nil {
			return
		}

		// Return the result of _parseDefault.
		err = _parseDefault(filepath.Join(projectPath, "config", file+".json"), cfg)
		if err != nil {
			return
		}
	})

	// Return nil to indicate the operation completed successfully.
	return err
}
