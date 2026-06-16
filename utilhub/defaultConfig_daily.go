package utilhub

import "sync"

// 🧪 Centralized feature toggle configuration for B Plus Tree unit testing.
var (
	// 🧪 Create a config instance for B plus tree unit testing and parse toggle switch settings. (开关设定值)
	_toggleConfig     = BptreeToggleConfig{}
	_onesToggleConfig sync.Once // Prevent configuration from being overwritten.
	// _toggleParseErr stores any error returned.
	_toggleParseErr = ParseToggle(&_toggleConfig)
)

// 🧪 Default configuration settings for daily automated B Plus Tree testing.
var (
	// 🧪 Create a config instance for B plus tree unit testing and parse default values. (单元测试设定值)
	_unitTestConfig     = BptreeUnitTestConfig{}
	_onesUnitTestConfig sync.Once // Prevent configuration from being overwritten.
	// _configParseErr stores any error returned.
	_configParseErr = ParseDefaultManual(&_unitTestConfig)
)

// 🧪 Initialize default test parameters.
func init() {
	if _configParseErr != nil {
		panic(_configParseErr)
	}
}

func ForceReloadConfig() {
	_unitTestConfig = BptreeUnitTestConfig{}
	_configParseErr = ParseDefault(&_unitTestConfig)
}

func GetDefaultConfig() BptreeUnitTestConfig {
	_configParseErr = ParseDefaultManual(&_unitTestConfig)

	if _configParseErr != nil {
		panic(_configParseErr)
	}

	return _unitTestConfig
}

func SetRandomTotalCount(value int64) {
	_unitTestConfig.Parameters.RandomTotalCount = value
}

func GetRandomTotalCount() int64 {
	return _unitTestConfig.Parameters.RandomTotalCount
}
