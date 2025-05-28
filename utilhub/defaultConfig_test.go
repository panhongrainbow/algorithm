package utilhub

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
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

func Test_SetFieldValue(t *testing.T) {
	t.Run("explain the function of SetFieldValue", func(t *testing.T) {
		var modified int16
		v := reflect.ValueOf(&modified).Elem()
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		require.Equal(t, modified, int16(100))
	})

	t.Run("explain how to test SetFieldValue", func(t *testing.T) {
		var modified int16
		var container interface{} = modified // 测试时，会传送 interface{}
		fieldType := reflect.TypeOf(container)
		v := reflect.New(fieldType).Elem()
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		modified = v.Interface().(int16)
		require.Equal(t, modified, int16(100))
	})

	t.Run("explain how to compare reflect.Value and interface{}", func(t *testing.T) {
		var modified int16
		var container interface{} = modified
		fieldType := reflect.TypeOf(container)
		v := reflect.New(fieldType).Elem()
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		container = v.Interface()
		require.Equal(t, container, interface{}(int16(100)))
	})

	type testCase struct {
		name      string
		fieldType interface{}
		input     string
		expected  interface{}
		shouldErr bool
	}

	tests := []testCase{
		{"int 100", 0, "100", 100, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldType := reflect.TypeOf(tc.fieldType)
			v := reflect.New(fieldType).Elem()
			// v := reflect.ValueOf(&tc.fieldType).Elem() // 不能这样写
			err := setFieldValue(v, tc.input)
			if tc.shouldErr {
				require.Error(t, err)
				return
			}
			tc.fieldType = v.Interface()
			require.NoError(t, err)
			require.Equal(t, tc.expected, tc.fieldType)
		})
	}
}
