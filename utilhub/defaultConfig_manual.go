package utilhub

func GetDefaultConfig() TestAutoConfig {
	_autoParseErr = ParseDefaultManual(&_autoConfig)

	if _autoParseErr != nil {
		panic(_autoParseErr)
	}

	return _autoConfig
}
