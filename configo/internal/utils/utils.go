package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

func GetNameOfTheObject(obj any) string {
	configNames := strings.Split(reflect.TypeOf(obj).String(), ".")
	return configNames[len(configNames)-1]
}

func IsValidJsonString(obj any) bool {
	str, ok := obj.(string)
	if !ok {
		return false
	}
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}
