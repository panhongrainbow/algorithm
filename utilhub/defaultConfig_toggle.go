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

// 🧪 Centralized feature toggle configuration for B Plus Tree unit testing. (测试模式总开关)
var (
	// 🧪 Create a config instance for B plus tree unit testing and parse toggle switch settings.
	_toggleConfig     = TestToggleConfig{}
	_onesToggleConfig sync.Once // Prevent configuration from being overwritten.
	// _toggleParseErr stores any error returned.
	_toggleParseErr = ParseToggle(&_toggleConfig)
)

// 🧪 Default configuration settings for daily automated B Plus Tree testing. (自动测试模式)
var (
	// 🧪 Create a config instance for B plus tree automatic unit testing and parse default values.
	_autoConfig     = TestAutoConfig{}
	_onesAutoConfig sync.Once // Prevent configuration from being overwritten.
	// _autoParseErr stores any error returned.
	_autoParseErr = ParseDefaultManual(&_autoConfig)
)

// 🧪 Initialize default test parameters.
func init() {
	if _autoParseErr != nil {
		fmt.Printf("autoParseErr: %v\n", _autoParseErr)
		panic(_autoParseErr)
	}
}

// ParseDefaultManual may load either default configuration values or manually specified configuration values, depending on Toggle Config.
func ParseDefaultManual(cfg AutoConfig) error {
	switch testMech {
	case "auto":
		return ParseAuto(cfg)
	case "manual":
		for _, manualConfig := range _manualTestConfig {
			fmt.Println(manualConfig.Record.ManualRecordDate, manualConfig.Record.ManualRecordFile)
			if manualConfig.Record.ManualRecordDate == testDate &&
				manualConfig.Record.ManualRecordFile == testFile {
				_autoConfig = manualConfig
				return nil
			}
		}
	}

	return nil
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
