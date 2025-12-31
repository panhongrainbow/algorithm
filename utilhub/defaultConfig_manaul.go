package utilhub

import (
	"fmt"
	"path/filepath"
)

var (
	// 🧪 Create a config instance for B plus tree unit testing and parse default values.
	_manualTestConfig = BptreeUnitTestConfig{}
	_manualParseErr   = ParseManual(&_manualTestConfig)
)

type Scenario20251220 interface{}

// ParseManual ⛏️ loads 载入之前的错误案例
func ParseManual(cfg Scenario20251220) error {
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

	fmt.Println(projectPath, file)

	// Return the result of _parseDefault.
	err = _parseDefault(filepath.Join(projectPath, "config", file+".json"), cfg)
	if err != nil {
		return err
	}

	// If the record is configured to be inside the project directory,
	// prepend the project path to the test record path
	if cfg.(*BptreeUnitTestConfig).Record.IsInsideProject == true {
		cfg.(*BptreeUnitTestConfig).Record.TestRecordPath = filepath.Join(projectPath, cfg.(*BptreeUnitTestConfig).Record.TestRecordPath)
	}

	fmt.Println(cfg)

	// Return nil to indicate the operation completed successfully.
	return nil

}
