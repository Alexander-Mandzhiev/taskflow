package txmanager

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestIsSerializationError(t *testing.T) {
	t.Run("deadlock", func(t *testing.T) {
		err := &mysql.MySQLError{Number: 1213, Message: "Deadlock"}
		if !isSerializationError(err) {
			t.Error("isSerializationError(1213) = false, want true")
		}
	})
	t.Run("lock wait timeout", func(t *testing.T) {
		err := &mysql.MySQLError{Number: 1205, Message: "Lock wait timeout"}
		if !isSerializationError(err) {
			t.Error("isSerializationError(1205) = false, want true")
		}
	})
	t.Run("other mysql error", func(t *testing.T) {
		err := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
		if isSerializationError(err) {
			t.Error("isSerializationError(1062) = true, want false")
		}
	})
	t.Run("non-mysql error", func(t *testing.T) {
		err := errors.New("generic error")
		if isSerializationError(err) {
			t.Error("isSerializationError(generic) = true, want false")
		}
	})
	t.Run("nil", func(t *testing.T) {
		if isSerializationError(nil) {
			t.Error("isSerializationError(nil) = true, want false")
		}
	})
}
