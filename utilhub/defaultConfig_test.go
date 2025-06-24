package utilhub

import (
	"fmt"
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

func Test_DefaultConfig(t *testing.T) {
	cfg := &BptreeUnitTestConfig{}

	// Get the default configuration directory.
	path, err := GetProjectDir(filepath.Join(ProjectName, "/utilhub/test_defaultConfig"))
	require.NoError(t, err)

	// Get the struct name to use as the filename.
	file, err := GetDefaultStructName(cfg)
	require.NoError(t, err)

	err = _defaultConfig2file(cfg, filepath.Join(path, file+".json"), true)
	require.NoError(t, err)

	return

	err = _ParseDefault(path+"/"+file+".json", cfg)
	require.NoError(t, err)

	if err := ParseDefault(cfg); err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
