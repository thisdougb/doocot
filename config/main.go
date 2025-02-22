// Package config provides a simple way to manage configuration
// parameters. Precendence is envrionment overrides config-file
// overrides default.
package config

import (
	"os"
	"strconv"
)

var defaultValues = map[string]interface{}{
	"DEBUG":       true,
	"DOOCOT_HOST": "http://localhost:8080",
}

func StringValue(key string) string {
	if defaultValue, ok := defaultValues[key]; ok {
		return getEnvVar(key, defaultValue.(string)).(string)
	}
	return ""
}

// ValueAsInt gets a string value from the env or default
func Int64Value(key string) int64 {

	if defaultValue, ok := defaultValues[key]; ok {
		return getEnvVar(key, defaultValue.(int64)).(int64)
	}
	return 0
}

// ValueAsInt gets a string value from the env or default
func Int32Value(key string) int32 {

	if defaultValue, ok := defaultValues[key]; ok {
		return getEnvVar(key, defaultValue.(int32)).(int32)
	}
	return 0
}

// ValueAsInt gets a string value from the env or default
func IntValue(key string) int {

	if defaultValue, ok := defaultValues[key]; ok {
		return getEnvVar(key, defaultValue.(int)).(int)
	}
	return 0
}

// ValueAsBool gets a string value from the env or default
func BoolValue(key string) bool {

	if defaultValue, ok := defaultValues[key]; ok {
		return getEnvVar(key, defaultValue.(bool)).(bool)
	}
	return false
}

func getEnvVar(key string, fallback interface{}) interface{} {

	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	switch fallback.(type) {
	case string:
		return value
	case bool:
		valueAsBool, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return valueAsBool
	case int:
		valueAsInt, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return valueAsInt
	}
	return fallback
}
