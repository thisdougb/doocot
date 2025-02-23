//go:build dev

package config

import (
	"os"
	"testing"
)

// adds our test values, when this file is included with go -tags
func init() {
	defaultValues["_TEST_INT_VALUE"] = 10
	defaultValues["_TEST_STR_VALUE"] = "AAA"
	defaultValues["_TEST_BOOL_VALUE"] = false
}

// Test string value type
func TestString(t *testing.T) {

	// test: unset potential env var, this should return the Str value in defaultValue[map]
	os.Unsetenv("_TEST_STR_VALUE")
	if "AAA" != StringValue("_TEST_STR_VALUE") {
		t.Errorf("AAA does not match %v", StringValue("_TEST_STR_VALUE"))
	}

	// test: now override the defaultValue[map] using an env var value
	os.Setenv("_TEST_STR_VALUE", "hello")
	if "hello" != StringValue("_TEST_STR_VALUE") {
		t.Errorf("hello does not match %v", StringValue("_TEST_STR_VALUE"))
	}
	os.Unsetenv("_TEST_STR_VALUE")
}

// Test int value type
func TestInt(t *testing.T) {

	// test: unset potential env var, this should return the int value in defaultValue[map]
	os.Unsetenv("_TEST_INT_VALUE")
	if 10 != IntValue("_TEST_INT_VALUE") {
		t.Errorf("10 does not match %v", IntValue("_TEST_INT_VALUE"))
	}

	// test: now override the defaultValue[map] using an env var value
	os.Setenv("_TEST_INT_VALUE", "20")
	if 20 != IntValue("_TEST_INT_VALUE") {
		t.Errorf("20 does not match %v", IntValue("_TEST_INT_VALUE"))
	}
	os.Unsetenv("_TEST_INT_VALUE")

	// test: now we use a non-int env var, which should be ignored
	os.Setenv("_TEST_INT_VALUE", ";")
	if 10 != IntValue("_TEST_INT_VALUE") {
		t.Errorf("10 does not match %v", IntValue("_TEST_INT_VALUE"))
	}
	os.Unsetenv("_TEST_INT_VALUE")

}

// Test bool value type
func TestBool(t *testing.T) {

	// test: unset potential env var, this should return the int value in defaultValue[map]
	os.Unsetenv("_TEST_BOOL_VALUE")
	if false != BoolValue("_TEST_BOOL_VALUE") {
		t.Errorf("false does not match %v", BoolValue("_TEST_BOOL_VALUE"))
	}

	// test: now override the defaultValue[map] using an env var value
	os.Setenv("_TEST_BOOL_VALUE", "true")
	if true != BoolValue("_TEST_BOOL_VALUE") {
		t.Errorf("true does not match %v", BoolValue("_TEST_BOOL_VALUE"))
	}
	os.Unsetenv("_TEST_BOOL_VALUE")

	// test: now we use a non-int env var, which should be ignored
	os.Setenv("_TEST_BOOL_VALUE", "hello")
	if false != BoolValue("_TEST_BOOL_VALUE") {
		t.Errorf("false does not match %v", BoolValue("_TEST_BOOL_VALUE"))
	}
	os.Unsetenv("_TEST_BOOL_VALUE")
}

func TestGetEnvVar(t *testing.T) {

	// test: when we set a str env var, should should get that value
	os.Setenv("_TEST_STR_NEW", "isset")
	if "isset" != getEnvVar("_TEST_STR_NEW", "isset") {
		t.Errorf("isset does not match %v", getEnvVar("_TEST_STR_NEW", "isset"))
	}
	os.Unsetenv("_TEST_STR_NEW")

	// test: when no env var exists we should use the fallback value in 2nd arg
	if "fallback" != getEnvVar("_TEST_STR_NEW", "fallback") {
		t.Errorf("fallback does not match %v", getEnvVar("_TEST_STR_NEW", "fallback"))
	}

	// test: when we set an int env var, should should get that value
	os.Setenv("_TEST_INT_NEW", "32")
	if 32 != getEnvVar("_TEST_INT_NEW", 1) {
		t.Errorf("32 does not match %v", getEnvVar("_TEST_INT_NEW", 1))
	}
	os.Unsetenv("_TEST_INT_NEW")

	// test: when we set an bool env var, should should get that value
	os.Setenv("_TEST_BOOL_NEW", "true")
	if true != getEnvVar("_TEST_BOOL_NEW", false) {
		t.Errorf("true does not match %v", getEnvVar("_TEST_BOOL_NEW", false))
	}
	os.Unsetenv("_TEST_BOOL_NEW")

	// test: when we set a non-int env var, should should get the fallback value
	// this is because we don't convert non-int automatically in getEnvVar
	os.Setenv("TEST_UNKNOWN", "2.2")
	if 1.1 != getEnvVar("TEST_UNKNOWN", 1.1) {
		t.Errorf("1.1 does not match %v", getEnvVar("TEST_UNKNOWN", 1.1))
	}
	os.Unsetenv("TEST_UNKNOWN")
}
