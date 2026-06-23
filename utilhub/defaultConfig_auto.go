package utilhub

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

// ParseAuto ⛏️ loads the default configuration from struct tags and applies it to the provided struct.
func ParseAuto(cfg AutoConfig) error {
	// Prepare the variable outside of the closure function.
	var err error
	var projectPath, file string

	// Use Golang's sync.Once to prevent the setting from being overwritten.
	_onesAutoConfig.Do(func() {

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
		err = _parseAuto(filepath.Join(projectPath, "config", file+".json"), cfg)
		if err != nil {
			return
		}

		// If the record is configured to be inside the project directory,
		// prepend the project path to the test record path
		if cfg.(*TestAutoConfig).Record.IsInsideProject == true {
			cfg.(*TestAutoConfig).Record.TestRecordPath = filepath.Join(projectPath, cfg.(*TestAutoConfig).Record.TestRecordPath)
		}

		// When the mechanism is set to automatic, the test subdirectory will automatically be set to the current date.
		if cfg.(*TestAutoConfig).Mechanism == "auto" {
			cfg.(*TestAutoConfig).Record.ManualRecordDate = TestTimeString("2006-01-02", "Asia/Shanghai")
		}

		// HACK => Below is the test code.
		cfg.(*TestAutoConfig).Parameters.BpWidth = []int{3, 10, 11, 12, 13}
	})

	// Return nil to indicate the operation completed successfully.
	return err
}

// _parseAuto ⛏️ loads the default configuration from struct tags and applies it to the provided struct.
// Configuration from the file, if the file exists, and applies and overwrites the struct. (以文件的配置为主,结构体配置为次)
func _parseAuto(filePath string, cfg AutoConfig) error {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	// Read the default configuration from the file.
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Do nothing; the tag will be handled later by applyDefaults. (之后由 applyDefaults 取 tag 决定)
		}
		return err
	}

	// Unmarshal the JSON data into the provided config and overwrite the default values.
	if err := json.Unmarshal(file, cfg); err != nil {
		return err
	}

	// [applyDefaults] applies the default values from struct tags to the provided config. (主要逻辑)
	if err := applyDefaults(cfg); err != nil {
		return err
	}

	// No error occurred, return nil.
	return nil
}
