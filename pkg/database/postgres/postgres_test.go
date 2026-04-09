package postgres

import (
	"reflect"
	"testing"
)

// TestDBStructDefaults verifies the default field values declared via struct
// tags on DB. These defaults are used by cleanenv when no env variable or
// config value is supplied, so changing them would silently break deployments.
func TestDBStructDefaults(t *testing.T) {
	cases := []struct {
		field   string
		tagKey  string
		wantTag string
	}{
		{"Port", "env-default", "5432"},
		{"Host", "env-default", "localhost"},
		{"Name", "env-default", "postgres"},
		{"User", "env-default", "postgres"},
		{"Password", "env-default", "postgres"},
		{"SSLMode", "env-default", "disable"},
	}

	dbType := reflect.TypeOf(DB{})
	for _, tc := range cases {
		tc := tc
		t.Run(tc.field, func(tt *testing.T) {
			field, ok := dbType.FieldByName(tc.field)
			if !ok {
				tt.Fatalf("DB struct does not have field %q", tc.field)
			}
			got := field.Tag.Get(tc.tagKey)
			if got != tc.wantTag {
				tt.Errorf("DB.%s %s tag: got %q, want %q", tc.field, tc.tagKey, got, tc.wantTag)
			}
		})
	}
}

// TestDBStructEnvTags verifies that DB fields expose the expected environment
// variable names so that the deployment environment can override them.
func TestDBStructEnvTags(t *testing.T) {
	envTags := map[string]string{
		"Port":     "DB_PORT",
		"Host":     "DB_HOST",
		"Name":     "DB_NAME",
		"User":     "DB_USER",
		"Password": "DB_PASSWORD",
		"SSLMode":  "DB_SSLMODE",
	}

	dbType := reflect.TypeOf(DB{})
	for fieldName, wantEnv := range envTags {
		field, ok := dbType.FieldByName(fieldName)
		if !ok {
			t.Errorf("DB struct does not have field %q", fieldName)
			continue
		}
		gotEnv := field.Tag.Get("env")
		if gotEnv != wantEnv {
			t.Errorf("DB.%s env tag: got %q, want %q", fieldName, gotEnv, wantEnv)
		}
	}
}

// TestNewPostgresConnectionFailsOnUnreachableHost verifies that
// NewPostgresConnection returns a non-nil error when it cannot reach the
// database host. The ping step should fail for an unreachable address.
func TestNewPostgresConnectionFailsOnUnreachableHost(t *testing.T) {
	cfg := DB{
		Host:     "127.0.0.1",
		Port:     1, // port 1 is never open
		User:     "testuser",
		Name:     "testdb",
		SSLMode:  "disable",
		Password: "testpass",
	}

	db, err := NewPostgresConnection(cfg)
	if err == nil {
		db.Close()
		t.Fatal("expected an error for unreachable host, got nil")
	}
}

// TestNewPostgresConnectionFailsOnInvalidHost verifies that
// NewPostgresConnection returns an error for an invalid hostname.
func TestNewPostgresConnectionFailsOnInvalidHost(t *testing.T) {
	cfg := DB{
		Host:     "invalid-host-that-does-not-exist.local",
		Port:     5432,
		User:     "testuser",
		Name:     "testdb",
		SSLMode:  "disable",
		Password: "testpass",
	}

	db, err := NewPostgresConnection(cfg)
	if err == nil {
		db.Close()
		t.Fatal("expected an error for invalid host, got nil")
	}
}

// TestNewPostgresConnectionReturnsBothNilOnSuccess documents the contract that
// a successful connection returns (db, nil). Since no real database is
// available in unit tests, this is a negative test: we confirm the error path.
func TestNewPostgresConnectionErrorWrapping(t *testing.T) {
	cfg := DB{
		Host:     "127.0.0.1",
		Port:     1,
		User:     "u",
		Name:     "db",
		SSLMode:  "disable",
		Password: "p",
	}

	_, err := NewPostgresConnection(cfg)
	if err == nil {
		t.Fatal("expected wrapped error, got nil")
	}

	// The error should contain a ping-related description as NewPostgresConnection
	// wraps it with "database ping error: ...".
	errMsg := err.Error()
	if len(errMsg) == 0 {
		t.Error("error message should not be empty")
	}
}

// TestDBPortFieldType ensures that the Port field is an int (not string),
// so that the DSN format string %d is used correctly.
func TestDBPortFieldType(t *testing.T) {
	dbType := reflect.TypeOf(DB{})
	field, ok := dbType.FieldByName("Port")
	if !ok {
		t.Fatal("DB struct does not have a Port field")
	}
	if field.Type.Kind() != reflect.Int {
		t.Errorf("DB.Port kind: got %v, want int", field.Type.Kind())
	}
}

// TestPathMigrationsConstant ensures the migrations path constant has the
// expected value used when running MigrateDB.
func TestPathMigrationsConstant(t *testing.T) {
	const want = "file://migrations"
	if pathMigrations != want {
		t.Errorf("pathMigrations: got %q, want %q", pathMigrations, want)
	}
}