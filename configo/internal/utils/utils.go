package utils

import (
	"reflect"
	"strings"
)

func GetNameOfTheObject(obj any) string {
	configNames := strings.Split(reflect.TypeOf(obj).String(), ".")
	return configNames[len(configNames)-1]
}
