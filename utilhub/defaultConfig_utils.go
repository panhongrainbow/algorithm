package utilhub

// ForceReloadAutoTest ⛏️ forces the configuration to be reloaded for testing.
func ForceReloadAutoTest() {
	_autoConfig = TestProcessConfigType{}
	_autoParseErr = ParseAuto(&_autoConfig)
}

func SetAutoRandomTotalCount(value int64) {
	_autoConfig.Parameters.RandomTotalCount = value
}

func GetAutoRandomTotalCount() int64 {
	return _autoConfig.Parameters.RandomTotalCount
}
