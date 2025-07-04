package utilhub

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"reflect"
	"testing"
)

// Test_SetFieldValue tests the setFieldValue function.
// setFieldValue is a helper function to set a field value using reflection.
func Test_SetFieldValue(t *testing.T) {
	// The first part demonstrates usage examples, the second part contains actual test cases.
	// (先测试，再说明)

	// Explanation 1: Modifying existing values with reflect.Value.
	// (很简单的原则，用 ValueOf 修改数值，接近底层)
	t.Run("demonstrates modifying existing values", func(t *testing.T) {
		// When you need to modify an existing value, use ValueOf to get a settable reflected value.

		var modified int16
		v := reflect.ValueOf(&modified).Elem() // use ValueOf to set value.
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		require.Equal(t, modified, int16(100))
	})

	// Explanation 2: Creating new values from type information.
	// (用 TypeOf 建立新值，TypeOf 可以获取最完整类型信息)
	t.Run("demonstrates creating new values from types", func(t *testing.T) {
		// When you need to create a new value, use TypeOf to get type information, then create via reflect.New.

		var modified int16
		var container interface{} = modified
		fieldType := reflect.TypeOf(container) // use TypeOf to create new value.
		v := reflect.New(fieldType).Elem()
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		modified = v.Interface().(int16)
		require.Equal(t, modified, int16(100))
	})

	// Explanation 3: Working with interface{} conversions.
	t.Run("demonstrates interface{} conversion", func(t *testing.T) {
		// Finally, tests comparison between two interface{} values and verifies the result.

		var modified int16
		var container interface{} = modified
		fieldType := reflect.TypeOf(container)
		v := reflect.New(fieldType).Elem()
		err := setFieldValue(v, "100")
		require.NoError(t, err)
		container = v.Interface() // Convert Value back to interface{}
		require.Equal(t, container, interface{}(int16(100)))
	})

	// Comprehensive test cases.
	type testCase struct {
		name      string
		fieldType interface{}
		input     string
		expected  interface{}
		shouldErr bool
	}

	tests := []testCase{
		// Boolean tests
		{"bool true", false, "true", true, false},

		// Signed integer tests
		{"int16 100", int16(0), "100", int16(100), false},
		{"int64 -100", int64(0), "-100", int64(-100), false},

		// Unsigned integer tests
		{"uint32 100", uint32(0), "100", uint32(100), false},
		{"uint32 -100 (invalid)", uint32(0), "-100", uint32(0), true},

		// Floating point tests
		{"float32 100.001", float32(0), "100.001", float32(100.001), false},
		{"float64 -100.001", float64(0), "-100.001", -100.001, false},

		// Complex number (expected to fail)
		{"complex 3+4i (unsupported)", complex(3, 4), "3+4i", complex(0, 0), true},

		// Array tests
		{"int array [1,2,3,4,5]", [5]int{}, "1,2,3,4,5", [5]int{1, 2, 3, 4, 5}, false},
		{"string array [a,b,c,d,e]", [5]string{}, "a,b,c,d,e", [5]string{"a", "b", "c", "d", "e"}, false},
	}

	// Test execution loop
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// The testing process below is similar to the one in
			// "Explanation 3: Working with interface{} conversions"
			// (这里和说明三相似)

			fieldType := reflect.TypeOf(tc.fieldType)
			v := reflect.New(fieldType).Elem()
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

// Test_DefaultConfig validates configuration loading logic by:
// (1) Generating default configs to JSON,
// (2) Verifying file parsing into structs,
// (3) Checking field-level defaults, and
// (4) Testing value overrides. The mock _parseDefault isolates file operations,
// while ParseDefault's production-ready version would handle advanced binding.
func Test_DefaultConfig(t *testing.T) {
	cfg := &testConfig{}

	// Retrieve the default configuration directory path.
	defaultPath, err := GetProjectDir(filepath.Join(ProjectName, "/temp/test_defaultConfig"))
	require.NoError(t, err)

	// Retrieve the modified configuration directory path.
	modifiedPath, err := GetProjectDir(filepath.Join(ProjectName, "/temp/test_modified_defaultConfig"))
	require.NoError(t, err)

	// Get the struct name to be used as the filename.
	fileName, err := GetDefaultStructName(cfg)
	require.NoError(t, err)

	// Save the default configuration to a JSON file in the specified test directory.
	err = _defaultConfig2file(cfg, filepath.Join(defaultPath, fileName+".json"), true)
	require.NoError(t, err)

	// Initialize a FileNode instance for directory operations.
	fm := FileNode{}
	// Create the default directory if it doesn't exist.
	err = fm.MkDir(defaultPath).Error()
	// Navigate to the default directory.
	fm = fm.Goto(defaultPath)
	require.NoError(t, err)
	// List all files in the default directory.
	_, fileList, err := fm.List()
	require.NoError(t, err)
	// Verify only the expected JSON file exists in the directory.
	require.Equal(t, []string{fileName + ".json"}, fileList)

	/*
		_parseDefault is an internal mock version of ParseDefault, used exclusively in test environments.
		The ParseDefault function provides a convenient way to automatically associate a struct with its configuration file and
		load default values upon initialization.

		_parseDefault 是 ParseDefault 在测试环境下的模拟实现，主要用于单元测试隔离。
		ParseDefault 作为配置加载的核心方法，支持声明式配置绑定 - 只需传入结构体实例，即可自动关联同名配置文件并完成默认值注入。
	*/

	// Parse the default configuration from the JSON file into the struct.
	err = _parseDefault(defaultPath+"/"+fileName+".json", cfg) //
	require.NoError(t, err)

	/*
		This comment block shows the expected JSON structure.
		The configuration includes server, database, and features sections.
		Each section contains specific fields with default values.
	*/

	/*
		{
		  "server": {
		    "host": "localhost",
		    "port": 8080
		  },
		  "database": {
		    "url": "postgres1://localhost:5432/mydb",
		    "username": "admin",
		    "password": "password",
		    "pool_size": 10
		  },
		  "features": [
		    "feature1",
		    "feature2",
		    "feature3"
		  ]
		}
	*/

	// Reparse the same configuration file to ensure consistency.
	err = _parseDefault(defaultPath+"/"+fileName+".json", cfg)
	require.NoError(t, err)

	// Validate all server-related fields match expected defaults.
	require.Equal(t, "localhost", cfg.Server.Host)
	require.Equal(t, 8080, cfg.Server.Port)
	// Validate all database-related fields match expected defaults.
	require.Equal(t, "postgres1://localhost:5432/mydb", cfg.Database.URL)
	require.Equal(t, "admin", cfg.Database.Username)
	require.Equal(t, "password", cfg.Database.Password)
	require.Equal(t, 10, cfg.Database.PoolSize)
	// Validate the feature list matches expected defaults.
	require.Equal(t, []string{"feature1", "feature2", "feature3"}, cfg.Features)

	/*
		This comment block shows the modified JSON structure.
		The modified configuration contains different values for most fields.
		Note the pool_size field is intentionally omitted to use the default value.
	*/

	/*
		{
		  "server": {
		    "host": "product",
		    "port": 8085
		  },
		  "database": {
		    "url": "postgres2://localhost:5432/newdb",
		    "username": "user",
		    "password": "12345"
			"pool_size": 20 // Remove this field in the modified json file to use the default value.
		  },
		  "features": [
		    "feature5",
		    "feature6",
		    "feature7"
		  ]
		}
	*/

	// Parse the modified configuration from its JSON file.
	err = _parseDefault(modifiedPath+"/"+fileName+".json", cfg)
	require.NoError(t, err)

	// Validate server fields now reflect the modified values.
	require.Equal(t, "product", cfg.Server.Host)
	require.Equal(t, 8085, cfg.Server.Port)
	// Validate database fields reflect modifications while pool size remains default.
	require.Equal(t, "postgres2://localhost:5432/newdb", cfg.Database.URL)
	require.Equal(t, "user", cfg.Database.Username)
	require.Equal(t, "12345", cfg.Database.Password)
	require.Equal(t, 10, cfg.Database.PoolSize)
	// Validate the feature list now shows modified values.
	require.Equal(t, []string{"feature5", "feature6", "feature7"}, cfg.Features)

	// Clean up by removing the test JSON file.
	err = fm.RemoveFile(defaultPath, fileName+".json")
	require.NoError(t, err)
	// Verify the directory is now empty.
	_, fileList, err = fm.List()
	require.NoError(t, err)
	require.Equal(t, []string(nil), fileList)
}
