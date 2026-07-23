package utilhub

import (
	"fmt"
	"path/filepath"
	"sync"
)

// =====================================================================================================================
//	🛠️ Default Config (Tool)
// Default Config is a tool that tags struct fields with default values.
// (DefaultConfig是一个工具,用于标记结构体字段的默认值)
// =====================================================================================================================

// HACK => 暂时的设定
var (
	testMech = "auto"
	testDate = "2026-01-11"
	testFile = "SingleNodeEndurance.do_not_open"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> Managed Global Variables

// 🧪 Centralized feature toggle configuration for B Plus Tree unit testing. (测试模式总开关)
var (
	// 🧪 Create a config instance for B plus tree unit testing and parse toggle switch settings.
	_toggleConfig     = ToggleConfigType{}
	_onesToggleConfig sync.Once // Prevent configuration from being overwritten.
	// _toggleParseErr stores any error returned.
	_toggleParseErr = ParseToggle(&_toggleConfig)
)

// 🧪 Default configuration settings for daily automated B Plus Tree testing. (自动测试模式)
var (
	// 🧪 Create a config instance for B plus tree automatic unit testing and parse default values.
	_autoConfig     = TestProcessConfigType{}
	_onesAutoConfig sync.Once // Prevent configuration from being overwritten.
	// _autoParseErr stores any error returned.
	_autoParseErr = ParseAuto(&_autoConfig)
)

// 🧪 Initialize default test parameters.
func init() {
	if _autoParseErr != nil {
		fmt.Printf("autoParseErr: %v\n", _autoParseErr)
		panic(_autoParseErr)
	}
}

// 🧪 Default configuration settings for manual B Plus Tree testing. (《手动测试模式)
var (
	// 🧪 Create a config instance for B plus tree unit testing and collect many previous failure scenarios.
	_manualConfig     []ManualConfigType // This time, using a slice that gathers various different failure scenarios.
	_onesManualConfig sync.Once          // Prevent configuration from being overwritten.
	// _manualParseErr stores any error returned.
	_manualParseErr = ParseManual(&_manualConfig)
)

// 🧪 Initialize manual test parameters.
func init() {
	if _manualParseErr != nil {
		panic(_manualParseErr)
	}
}

// >>>>> >>>>> >>>>> >>>>> >>>>> Toggle Logic (总)

// ToggleConfig represents the master switch for configuration.
// It determines whether the system should use the default configuration or a manually provided configuration. (开关配置)
type ToggleConfig any

// ToggleConfigType ⛏️ defines the toggle settings for Bptree tests.
type ToggleConfigType struct {
	Mechanism string `json:"mechanism" default:"auto"` // 🧪 When the mechanism selection is set to `auto`, all tests will be conducted.
}

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
		err = _parseAuto(filepath.Join(projectPath, "config", file+".json"), cfg)
		if err != nil {
			return
		}
	})

	// Return nil to indicate the operation completed successfully.
	return err
}

// ManualConfig is the instance, which refers to the file name under the config directory.
// type ManualConfig interface{}

// ParseManual ⛏️ loads previous failure scenarios.
/*
func ParseManual(cfg ManualConfig) error {
	return _parseManual(cfg)
}
*/

// ManualConfig ⛏️ is a type constraint that allows struct types to store default configuration values. (预设配置)
// type ManualConfig interface{}

// ParseDefaultManual may load either default configuration values or manually specified configuration values, depending on Toggle Config.
/*
func ParseDefaultManual(cfg AutoConfig) error {
	switch testMech {
	case "auto":
		return ParseAuto(cfg)
	case "manual":
		for _, manualConfig := range _manualConfig {
			fmt.Println(manualConfig.Record.ManualRecordDate, manualConfig.Record.ManualRecordFile)
			if manualConfig.Record.ManualRecordDate == testDate &&
				manualConfig.Record.ManualRecordFile == testFile {
				// _autoConfig = manualConfig
				return nil
			}
		}
	}

	return nil
}
*/
