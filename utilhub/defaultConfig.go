package utilhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type DefaultConfig interface{}

func Load(filePath string, cfg DefaultConfig) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	if err := applyDefaults(cfg); err != nil {
		return err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(file, cfg); err != nil {
		return err
	}

	return nil
}

func OverWrite(filePath string, cfg DefaultConfig) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	if err := applyDefaults(cfg); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

func loadDefaults(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := applyDefaults(field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		if err := setFieldValue(field, defaultTag); err != nil {
			return fmt.Errorf("field %s: %v", fieldType.Name, err)
		}
	}
	return nil
}

func applyDefaults(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := applyDefaults(field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		if err := setFieldValue(field, defaultTag); err != nil {
			return fmt.Errorf("field %s: %v", fieldType.Name, err)
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("cannot set field value")
	}

	switch field.Kind() {
	case reflect.Invalid:
		return errors.New("invalid field kind")
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Complex64, reflect.Complex128:
		return errors.New("unsupported field type: complex")
	case reflect.Array:
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
		return errors.New("unsupported field type: channel")
	case reflect.Func:
		return errors.New("unsupported field type: function")
	case reflect.Interface:
		return errors.New("unsupported field type: interface")
	case reflect.Map:
		return errors.New("unsupported field type: map")
	case reflect.Pointer:
		return errors.New("unsupported field type: pointer")
	case reflect.Struct:
		return errors.New("unsupported field type: struct")
	case reflect.UnsafePointer:
	case reflect.String:
		field.SetString(value)
	case reflect.Slice:
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
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}
