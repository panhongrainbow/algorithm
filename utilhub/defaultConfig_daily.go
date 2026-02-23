package utilhub

import "sync"

var (
	// 🧪 Create a config instance for B plus tree unit testing and parse default values.
	_unitTestConfig = BptreeUnitTestConfig{}
	_configParseErr = ParseDefaultManual(&_unitTestConfig)
	_ones           sync.Once // Prevent configuration from being overwritten.
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
