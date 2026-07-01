package utilhub

// ManualConfigType ⛏️ is a struct for manual test configuration.
type ManualConfigType struct {
	Mechanism string   `json:"mechanism" default:"auto"` // 🧪 When the mechanism selection is set to `auto`, all tests will be conducted.
	Record    struct { // 🧪 Record contains configurations related to test record storage.
		TestRecordPath   string `json:"testRecordPath" default:"/temp/test_record"` // 🧪 TestRecordPath specifies the directory path where test records will be saved.
		ManualRecordDate string `json:"manualRecordDate" default:"0000-00-00"`      // 🧪 Manual testing requires specifying the previous test date.
		ManualRecordFile string `json:"manualRecordFile" default:"empty"`           // 🧪 Manual testing requires specifying the previous test record file.
		IsInsideProject  bool   `json:"isInsideProject" default:"true"`             // 🧪 IsInsideProject indicates whether the test records are stored inside the project directory.
	} `json:"record"`
	Parameters struct { // Parameters contains configurations for test execution parameters.
		RandomTotalCount             int64 `json:"randomTotalCount" default:"7500000"`        // 🧪 RandomTotalCount represents the number of elements to be generated for random testing.
		RandomMin                    int64 `json:"randomMin" default:"10"`                    // 🧪 RandomMin represents the minimum value for generating random numbers.
		RandomHitCollisionPercentage int64 `json:"randomHitCollisionPercentage" default:"70"` // 🧪 Random number hit collision percentage.
		// Calculate the maximum random value.
		// randomTotalCount/randomHitCollisionPercentage*100 + randomMin = randomMax
		// 7500000 / 70 * 100 + 10 = 10714295
		RandomMax int64 `json:"randomMax" default:"10714295"` // 🧪 RandomMax represents the maximum value for generating random numbers.
		BpWidth   []int `json:"bpWidth" default:"3,4,5,6,7"`
	} `json:"parameters"`
	PoolStage struct { // This is primarily used to test boundary conditions.
		MinRemovals       int64 `json:"minRemovals" default:"5"`        // 🧪 Lower bound of items to remove in this stage.
		MaxRemovals       int64 `json:"maxRemovals" default:"50"`       // 🧪 Upper bound of items to remove in this stage.
		MinPreserveInPool int64 `json:"minPreserveInPool" default:"10"` // 🧪 Lower bound of items to remain in the pool after this stage.
		MaxPreserveInPool int64 `json:"maxPreserveInPool" default:"20"` // 🧪 Upper bound of items to remain in the pool after this stage.
	} `json:"poolStage"`
	CyclicStress struct { // metal fatigue style endurance test.
		CyclicStressCount int `json:"cyclicStressCount" default:"10"` // 🧪 Number of fatigue test cycles.
	} `json:"cyclicStress"`
	ManualTest struct { // 使用手动测试，重现之前的错误
		EnableBulkInsertDelete   bool `json:"enableBulkInsertDelete" default:"false"`
		EnableRandomizedBoundary bool `json:"enableRandomizedBoundary" default:"false"`
		EnableNodeEnduranceTest  bool `json:"enableNodeEnduranceTest" default:"false"`
	} `json:"manualTest"`
}

/*
func GetDefaultConfig() AutoConfigType {
	_autoParseErr = ParseDefaultManual(&_autoConfig)

	if _autoParseErr != nil {
		panic(_autoParseErr)
	}

	return _autoConfig
}
*/
