package validation

import (
	"testing"
)

// testStruct используется для проверки работы глобального Validator (теги как в DTO).
type testStruct struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Name     string `validate:"omitempty,max=255"`
}

func TestValidator_Struct(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := testStruct{
			Email:    "user@example.com",
			Password: "password123",
			Name:     "Alice",
		}
		err := Validator.Struct(v)
		if err != nil {
			t.Errorf("Validator.Struct(valid) = %v", err)
		}
	})

	t.Run("valid omitempty name", func(t *testing.T) {
		v := testStruct{
			Email:    "a@b.co",
			Password: "12345678",
			Name:     "",
		}
		err := Validator.Struct(v)
		if err != nil {
			t.Errorf("Validator.Struct(valid, empty name) = %v", err)
		}
	})

	t.Run("missing email", func(t *testing.T) {
		v := testStruct{
			Email:    "",
			Password: "password123",
		}
		err := Validator.Struct(v)
		if err == nil {
			t.Fatal("Validator.Struct(missing email) expected error")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		v := testStruct{
			Email:    "not-an-email",
			Password: "password123",
		}
		err := Validator.Struct(v)
		if err == nil {
			t.Fatal("Validator.Struct(invalid email) expected error")
		}
	})

	t.Run("short password", func(t *testing.T) {
		v := testStruct{
			Email:    "user@example.com",
			Password: "short",
		}
		err := Validator.Struct(v)
		if err == nil {
			t.Fatal("Validator.Struct(short password) expected error")
		}
	})

	t.Run("name too long", func(t *testing.T) {
		longName := ""
		for i := 0; i < 256; i++ {
			longName += "a"
		}
		v := testStruct{
			Email:    "user@example.com",
			Password: "password123",
			Name:     longName,
		}
		err := Validator.Struct(v)
		if err == nil {
			t.Fatal("Validator.Struct(name max=255) expected error")
		}
	})
}
