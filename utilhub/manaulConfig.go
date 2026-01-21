package utilhub

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

var (
	// 🧪 Create a config instance for B plus tree unit testing and collect many previous failure scenarios.
	_manualTestConfig []BptreeUnitTestConfig // This time, using a slice that gathers various different failure scenarios.
	_manualParseErr   = ParseManual(&_manualTestConfig)
)

// 🧪 Initialize manual test parameters.
func init() {
	if _manualParseErr != nil {
		panic(_manualParseErr)
	}
}

// ManualConfig is the instance, which refers to the file name under the config directory.
type ManualConfig interface{}

// ParseManual ⛏️ loads previous failure scenarios.
func ParseManual(cfg ManualConfig) error {
	return _parseManual(cfg)
}

// _parseManual ⛏️ loads the manual configurations from struct tags and applies it to the provided struct,
// prioritizing the specified values if available.
func _parseManual(cfg ManualConfig) error {
	// Get the default configuration directory.
	projectPath, err := GetProjectDir(filepath.Join(ProjectName))
	if err != nil {
		return err
	}

	// Get the struct name to use as the filename.
	file, err := GetDefaultStructName(&cfg)
	if err != nil {
		return err
	}

	// _applyDefaults2manualConfig will fill in the missing fields of the instance with default values.
	err = _applyDefaults2manualConfig(filepath.Join(projectPath, "config", file+".json"), cfg)
	if err != nil {
		return err
	}

	// If the record is configured to be inside the project directory,
	// prepend the project path to the test record path.
	arr := *cfg.(*[]BptreeUnitTestConfig)
	for i := 0; i < len(arr); i++ {
		if arr[i].Record.IsInsideProject == true {
			arr[i].Record.TestRecordPath = filepath.Join(projectPath, arr[i].Record.TestRecordPath)
		}
	}

	// Return nil to indicate the operation completed successfully.
	return nil
}

// _applyDefaults2manualConfig fills in default values for many previous failure scenarios.
func _applyDefaults2manualConfig(filePath string, cfg ManualConfig) error {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	// Read the default configuration from the file.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Fills in default values for the content.
	err = _applyDefaults2content(content, cfg)
	if err != nil {
		return err
	}

	// No error occurred, return nil.
	return nil
}

// _applyDefaults2content function fills in default values for the content of many previous failure scenarios.
func _applyDefaults2content(content []byte, cfg ManualConfig) error {
	// Unmarshal the JSON data into the provided config and overwrite the default values.
	if err := json.Unmarshal(content, cfg); err != nil {
		return err
	}

	// Extract each configuration file from the slice, and fill in the default values if any fields are missing.
	arr := cfg.(*[]BptreeUnitTestConfig)
	for i := 0; i < len(*arr); i++ {
		// [applyDefaults] applies the default values from struct tags to the provided config. (主要逻辑)
		if err := applyDefaults(&((*arr)[i])); err != nil {
			return err
		}
	}

	// No error occurred, return nil.
	return nil
}
