package utilhub

// =====================================================================================================================
//                  🛠️ Default Config Type (Tool)
// Default Config Type contains types for DefaultConfig, bptreeUnitTestConfig etc. (这里收集了 DefaultConfig 等类型)
// =====================================================================================================================

// DefaultConfig ⛏️ is a type constraint that allows struct types to store default configuration values. (预设配置)
type DefaultConfig interface{}

// BptreeUnitTestConfig ⛏️ is a struct for BpTree unit test configuration.
type BptreeUnitTestConfig struct {
	RandomTotalCount             int64 `json:"randomTotalCount" default:"7500000"`        // 🧪 randomTotalCount represents the number of elements to be generated for random testing.
	RandomMin                    int64 `json:"randomMin" default:"10"`                    // 🧪 randomMin represents the minimum value for generating random numbers.
	RandomHitCollisionPercentage int64 `json:"randomHitCollisionPercentage" default:"70"` // 🧪 random number hit collision percentage.
	// Calculate the maximum random value.
	// randomTotalCount/randomHitCollisionPercentage*100 + randomMin = randomMax
	// 7500000 / 70 * 100 + 10 = 10714295
	RandomMax int64 `json:"randomMax" default:"10714295"` // 🧪 randomMax represents the maximum value for generating random numbers.
}

// types for testing is as bellows: (以下是测试用的类型) ===== ===== ===== ===== ===== ===== ===== ===== =====

// testConfig ⛏️ is a test struct for DefaultConfig. (测试用的预设配置)
type testConfig struct {
	Server struct {
		Host string `json:"host" default:"localhost"`
		Port int    `json:"port" default:"8080"`
	} `json:"server"`
	Database struct {
		URL      string `json:"url" default:"postgres1://localhost:5432/mydb"`
		Username string `json:"username" default:"admin"`
		Password string `json:"password" default:"password"`
		PoolSize int    `json:"pool_size" default:"10"`
	} `json:"database"`
	Features []string `json:"features" default:"feature1,feature2,feature3"`
}
