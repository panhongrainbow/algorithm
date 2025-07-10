package utilhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// =====================================================================================================================
//	🛠️ Default Config (Tool)
// Default Config is a tool that tags struct fields with default values.
// (DefaultConfig是一个工具,用于标记结构体字段的默认值)
// =====================================================================================================================

// ParseDefault ⛏️ loads the default configuration from struct tags and applies it to the provided struct.
func ParseDefault(cfg DefaultConfig) error {
	// Get the default configuration directory.
	projectPath, err := GetProjectDir(filepath.Join(ProjectName))
	if err != nil {
		return err
	}

	// Get the struct name to use as the filename.
	file, err := GetDefaultStructName(&cfg)
	if err != nil {
		return err
	}

	// Return the result of _parseDefault.
	err = _parseDefault(filepath.Join(projectPath, "config", file+".json"), cfg)
	if err != nil {
		return err
	}

	// If the record is configured to be inside the project directory,
	// prepend the project path to the test record path
	if cfg.(*BptreeUnitTestConfig).Record.IsInsideProject == true {
		cfg.(*BptreeUnitTestConfig).Record.TestRecordPath = filepath.Join(projectPath, cfg.(*BptreeUnitTestConfig).Record.TestRecordPath)
	}

	// Return nil to indicate the operation completed successfully.
	return nil

}

// _parseDefault ⛏️ loads the default configuration from struct tags and applies it to the provided struct.
// Configuration from the file, if the file exists, and applies and overwrites the struct. (以文件的配置为主,结构体配置为次)
func _parseDefault(filePath string, cfg DefaultConfig) error {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	// Read the default configuration from the file.
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Do nothing; the tag will be handled later by applyDefaults. (之后由 applyDefaults 取 tag 决定)
		}
		return err
	}

	// Unmarshal the JSON data into the provided config and overwrite the default values.
	if err := json.Unmarshal(file, cfg); err != nil {
		return err
	}

	// [applyDefaults] applies the default values from struct tags to the provided config. (主要逻辑)
	if err := applyDefaults(cfg); err != nil {
		return err
	}

	// No error occurred, return nil.
	return nil
}

// applyDefaults ⛏️ applies the default values from struct tags to the provided config.
func applyDefaults(cfg interface{}) error {
	// Get the reflect.Value of the passed-in struct and dereference it.
	v := reflect.ValueOf(cfg).Elem()

	// Get the type information of the struct.
	t := v.Type()

	// Iterate through all fields in the struct.
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)     // This will be used later to get the actual value of the field. (在这里获取实际值)
		fieldType := t.Field(i) // This will be used later to get the default tag value. (在这里获取预设值)

		// If the field is a struct, recursively apply defaults to it.
		if field.Kind() == reflect.Struct {
			if err := applyDefaults(field.Addr().Interface()); err != nil { // (这里是递归)
				return err
			}
			continue
		}

		// Get the "default" tag value from the field.
		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		// Skip setting the value with defaultTag if it has already been loaded from the config file.
		// (如果之前已经从配置文件中读取到了值，就跳过，不再使用 defaultTag 设置)
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		// Set the field to the default value from the tag.
		if err := setFieldValue(field, defaultTag); err != nil {
			return fmt.Errorf("field %s: %v", fieldType.Name, err)
		}
	}

	// No error occurred, return nil.
	return nil
}

// setFieldValue ⛏️ sets the value of a field based on its type.
func setFieldValue(field reflect.Value, value string) error {
	// Return an error if the field cannot be set.
	if !field.CanSet() {
		return errors.New("cannot set field value")
	}

	// Determine the kind of the field and handle accordingly.
	switch field.Kind() {
	case reflect.Invalid:
		// Return an error if the field kind is invalid.
		return errors.New("invalid field kind")
	case reflect.Bool:
		// Parse and set a boolean value.
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Parse and set an integer value.
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Parse and set an unsigned integer value.
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		// Parse and set a floating-point number.
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Complex64, reflect.Complex128:
		// Complex numbers are not supported.
		return errors.New("unsupported field type: complex") // (先不管 复数)
	case reflect.Array:
		// Split the string into items and set each array element.
		items := strings.Split(value, ",")
		if len(items) != field.Len() {
			return fmt.Errorf("array length mismatch: expected %d, got %d", field.Len(), len(items))
		}
		for i := 0; i < field.Len(); i++ {
			item := strings.TrimSpace(items[i])
			if err := setFieldValue(field.Index(i), item); err != nil {
				return err
			}
		}
	case reflect.Chan:
		// Channels are not supported.
		return errors.New("unsupported field type: channel")
	case reflect.Func:
		// Functions are not supported.
		return errors.New("unsupported field type: function")
	case reflect.Interface:
		// Interfaces are not supported.
		return errors.New("unsupported field type: interface")
	case reflect.Map:
		// Maps are not supported.
		return errors.New("unsupported field type: map")
	case reflect.Pointer:
		// Pointers are not supported.
		return errors.New("unsupported field type: pointer")
	case reflect.Struct:
		// Structs are not supported.
		return errors.New("unsupported field type: struct")
	case reflect.UnsafePointer:
		// UnsafePointer are not supported.
		return errors.New("unsupported field type: UnsafePointer")
	case reflect.String:
		// Set a string value.
		field.SetString(value)
	case reflect.Slice:
		// Split the string into items and set each slice element.
		items := strings.Split(value, ",")
		slice := reflect.MakeSlice(field.Type(), len(items), len(items))
		for i, item := range items {
			item = strings.TrimSpace(item)
			if err := setFieldValue(slice.Index(i), item); err != nil {
				return err
			}
		}
		field.Set(slice)
	default:
		// Return an error for unsupported field types.
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	// Return nil to indicate the value was successfully set.
	return nil
}

// defaultConfig2file ⛏️ saves the default configuration to a JSON file.
func defaultConfig2file(cfg DefaultConfig, overwrite bool) error {
	// Get the default configuration directory.
	path, err := GetProjectDir("algorithm/config")
	if err != nil {
		return err
	}

	// Get the struct name to use as the filename.
	file, err := GetDefaultStructName(&cfg)
	if err != nil {
		return err
	}

	// Construct the full file path and save the configuration.
	return _defaultConfig2file(cfg, filepath.Join(path, file+".json"), overwrite)
}

// _defaultConfig2file ⛏️ overwrites the default configuration from struct tags and applies it to the specific file.
func _defaultConfig2file(cfg DefaultConfig, filePath string, overwrite bool) error {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	// Apply default values to any unset fields in the config.
	if err := applyDefaults(cfg); err != nil {
		return err
	}

	// Marshal the config into indented JSON format.
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// If overwrite is true, write the JSON data to the specified file path.
	if overwrite {
		// Write the JSON data to the specified file path.
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	}

	// Return nil to indicate the operation completed successfully.
	return nil
}

// GetDefaultStructName ⛏️ retrieves the name of the struct.
func GetDefaultStructName(cfg DefaultConfig) (string, error) {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return "", errors.New("config must be a pointer to a struct")
	}

	// Get the name of the struct.
	v := reflect.ValueOf(cfg).Elem()

	// Return the name of the struct.
	return v.Type().Name(), nil
}

// GetProjectDir ⛏️ retrieves the absolute path to the algorithm project's subdirectory.
func GetProjectDir(subDirs string) (string, error) {
	// Get the caller's file path (this file).
	_, callerPath, _, _ := runtime.Caller(0)

	// Convert to an absolute path.
	utilhubPath, err := filepath.Abs(callerPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Split path to find algorithm root and join with config directory.
	paths := strings.Split(utilhubPath, "algorithm")
	configPath := filepath.Join(paths[0], subDirs)

	// Return an absolute config path.
	return configPath, nil
}
