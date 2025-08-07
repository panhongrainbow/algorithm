package utilhub

var (
	// 🧪 Create a config struct for B plus tree unit testing and parse default values.
	_unitTestConfig = BptreeUnitTestConfig{}
	_configParseErr = ParseDefault(&_unitTestConfig)
)

func init() {
	if _configParseErr != nil {
		panic(_configParseErr)
	}
}

func GetDefaultConfig() BptreeUnitTestConfig {
	return _unitTestConfig
}
