package utilhub

import (
	"fmt"
	"testing"
)

type AppConfig struct {
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

func TestDefaultConfig(t *testing.T) {
	cfg := &AppConfig{}

	// 加载配置
	if err := Load("/home/panhong/go/src/github.com/panhongrainbow/algorithm/utilhub/default_config.json", cfg); err != nil {
		panic(err)
	}

	fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: %s (pool: %d)\n", cfg.Database.URL, cfg.Database.PoolSize)
	fmt.Printf("Features: %v\n", cfg.Features)
}
