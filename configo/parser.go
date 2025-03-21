package configo

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func parseTags(ctx context.Context, cfg any) ([]ConfigObject, error) {
	// Get the type and value of the given config
	cfgType := reflect.TypeOf(cfg)
	cfgValue := reflect.ValueOf(cfg)
	if cfgType.Kind() == reflect.Ptr {
		cfgType = cfgType.Elem()
		cfgValue = cfgValue.Elem()
	}

	// Ensure the provided config is a struct
	if cfgType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", cfgType.Kind())
	}

	var configs []ConfigObject
	var err error
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		value := cfgValue.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}
		configObj := ConfigObject{
			Name:        field.Tag.Get(string(CONFIG_TAG_NAME)),
			Type:        ConfigType(field.Tag.Get(string(CONFIG_TAG_TYPE))),
			Description: field.Tag.Get(string(CONFIG_TAG_DESCRIPTION)),
			Required:    field.Tag.Get(string(CONFIG_TAG_REQUIRED)) == "true",
		}

		if choices := field.Tag.Get(string(CONFIG_TAG_CHOICES)); choices != "" {
			configObj.Choices = parseChoices(choices)
		}

		if configObj.Type == CONFIG_TYPE_PARENT {
			configObj.Childrens, err = parseTags(ctx, value.Interface())
			if err != nil {
				return nil, err
			}
		} else {
			configObj.Value = value.Interface()
		}

		configs = append(configs, configObj)
	}

	return configs, nil
}

func parseChoices(choices string) []string {
	return strings.Split(choices, ",")
}

func marshal(ctx context.Context, cfg any) (string, error) {
	objects, err := parseTags(ctx, cfg)
	if err != nil {
		return "", err
	}
	for _, obj := range objects {
		if err := obj.Validate(); err != nil {
			return "", err
		}
	}
	bytes, err := json.Marshal(objects)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func unmarshal(ctx context.Context, value string, cfg any) error {
	// Parse the JSON into []ConfigObject
	var configObjects []ConfigObject
	if err := json.Unmarshal([]byte(value), &configObjects); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Populate the struct with values from configObjects
	return populateStruct(ctx, cfg, configObjects)
}

// populateStruct recursively sets struct field values based on parsed ConfigObject data.
func populateStruct(ctx context.Context, cfg any, configObjects []ConfigObject) error {

	if len(configObjects) == 0 {
		return nil
	}

	// Get the reflect type and value of cfg
	cfgType := reflect.TypeOf(cfg)
	cfgValue := reflect.ValueOf(cfg)

	// Ensure cfg is a pointer to a struct
	if cfgType.Kind() != reflect.Ptr || cfgType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got %s", cfgType.Kind())
	}

	// Get the actual struct type and value
	cfgType = cfgType.Elem()
	cfgValue = cfgValue.Elem()

	for _, configObj := range configObjects {
		// Find the field by its JSON tag name
		field, found := findFieldByTag(cfgType, CONFIG_TAG_NAME, configObj.Name)
		if !found {
			continue // Ignore fields that don't match
		}

		fieldValue := cfgValue.FieldByName(field.Name)

		// Ensure the field is settable
		if !fieldValue.CanSet() {
			continue
		}

		// If the field is a struct, recursively populate it
		if field.Type.Kind() == reflect.Struct {
			var fieldValuePtr any
			if fieldValue.Kind() == reflect.Ptr {
				fieldValuePtr = fieldValue.Interface()
			} else {
				fieldValuePtr = fieldValue.Addr().Interface()
			}
			err := populateStruct(ctx, fieldValuePtr, configObj.Childrens)
			if err != nil {
				return err
			}
		} else {
			// Set primitive field values based on type
			err := setFieldValue(fieldValue, configObj.Value)
			if err != nil {
				return fmt.Errorf("error setting field %s: %w", field.Name, err)
			}
		}
	}
	return nil
}

// findFieldByTag searches for a struct field by a given tag key and value.
func findFieldByTag(cfgType reflect.Type, tagKey ConfigTag, tagValue string) (reflect.StructField, bool) {

	tagKeyStr := string(tagKey)
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		if field.Tag.Get(tagKeyStr) == tagValue {
			return field, true
		}
	}
	return reflect.StructField{}, false
}

// setFieldValue sets a field value based on its type.
func setFieldValue(field reflect.Value, value any) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	// Convert value based on field type
	switch field.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		} else {
			return fmt.Errorf("expected string, got %T", value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := value.(float64); ok { // JSON numbers are decoded as float64
			field.SetInt(int64(v))
		} else {
			return fmt.Errorf("expected int, got %T", value)
		}
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		} else {
			return fmt.Errorf("expected bool, got %T", value)
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(field.Type(), 0, len(value.([]any)))
		for _, item := range value.([]any) {
			elem := reflect.New(field.Type().Elem()).Elem()
			if err := setFieldValue(elem, item); err != nil {
				return err
			}
			slice = reflect.Append(slice, elem)
		}
		field.Set(slice)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
