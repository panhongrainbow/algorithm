package utilhub

// =====================================================================================================================
//                  🛠️ Default Config Type (Tool)
// Default Config Type contains types for DefaultConfig, bptreeUnitTestConfig etc. (这里收集了 DefaultConfig 等类型)
// =====================================================================================================================

// types for testing is as bellows: (以下是测试用的类型) ===== ===== ===== ===== ===== ===== ===== ===== =====

// testConfig ⛏️ is a test struct for AutoConfig. (测试用的预设配置)
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
