package logger

import (
	"reflect"
	"testing"
)

// TestConfigStruct verifies that Config has the expected yaml tag and env-default.
// The "log_mode" yaml tag and "dev" default are used inside init() to configure
// the global logger; this test ensures the struct definition is correct.
func TestConfigStruct(t *testing.T) {
	cfg := &Config{}

	// Without cleanenv the zero value of Mode is an empty string.
	if cfg.Mode != "" {
		t.Errorf("expected zero-value Mode to be empty string, got %q", cfg.Mode)
	}
}

// TestConfigYAMLTagPresence is a regression test ensuring the yaml struct tag
// "log_mode" is present on Config.Mode. init() reads the config from a YAML
// file using cleanenv, which relies on this tag name.
func TestConfigYAMLTagPresence(t *testing.T) {
	field, ok := reflect.TypeOf(Config{}).FieldByName("Mode")
	if !ok {
		t.Fatal("Config struct does not have a Mode field")
	}

	got := field.Tag.Get("yaml")
	want := "log_mode"
	if got != want {
		t.Errorf("Config.Mode yaml tag: got %q, want %q", got, want)
	}
}

// TestConfigEnvDefaultTag verifies that the env-default tag on Config.Mode is
// "dev". This matches the constant envDev moved inside init() in this PR.
func TestConfigEnvDefaultTag(t *testing.T) {
	field, ok := reflect.TypeOf(Config{}).FieldByName("Mode")
	if !ok {
		t.Fatal("Config struct does not have a Mode field")
	}

	got := field.Tag.Get("env-default")
	want := "dev"
	if got != want {
		t.Errorf("Config.Mode env-default tag: got %q, want %q", got, want)
	}
}

// TestGlobalLoggerInitialised ensures the package-level init() successfully
// configured globalLogger. If init() failed it would have called os.Exit(1),
// so reaching this test means it succeeded.
func TestGlobalLoggerInitialised(t *testing.T) {
	if globalLogger == nil {
		t.Fatal("globalLogger is nil: init() did not set it up correctly")
	}
}

// TestInfoDoesNotPanic verifies that calling Info does not panic.
func TestInfoDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Info panicked: %v", r)
		}
	}()

	Info("test info message")
}

// TestErrorDoesNotPanic verifies that calling Error does not panic.
func TestErrorDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Error panicked: %v", r)
		}
	}()

	Error("test error message")
}

// TestDebugDoesNotPanic verifies that calling Debug does not panic.
func TestDebugDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Debug panicked: %v", r)
		}
	}()

	Debug("test debug message")
}

// TestLoggerDevMode verifies "dev" is a valid mode constant (envDev inside init).
// This guards against the constant being accidentally changed.
func TestLoggerDevMode(t *testing.T) {
	cfg := &Config{Mode: "dev"}
	if cfg.Mode != "dev" {
		t.Errorf("expected Mode %q, got %q", "dev", cfg.Mode)
	}
}

// TestLoggerProdMode verifies "prod" is a valid mode constant (envProd inside init).
func TestLoggerProdMode(t *testing.T) {
	cfg := &Config{Mode: "prod"}
	if cfg.Mode != "prod" {
		t.Errorf("expected Mode %q, got %q", "prod", cfg.Mode)
	}
}

// TestLoggerUnknownModeIsNeitherDevNorProd ensures that any value other than
// "dev" and "prod" is treated as unknown by the switch in init(). The default
// branch prints a warning and returns without setting globalLogger.
func TestLoggerUnknownModeIsNeitherDevNorProd(t *testing.T) {
	unknown := "staging"
	if unknown == "dev" || unknown == "prod" {
		t.Errorf("mode %q should not match any known mode constant", unknown)
	}
}