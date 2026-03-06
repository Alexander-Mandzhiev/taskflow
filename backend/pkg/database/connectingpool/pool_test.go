package connectingpool

import (
	"context"
	"testing"
)

func TestNew_InvalidDriver(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, "unknown_driver_xyz", "localhost:3306")
	if err == nil {
		t.Fatal("New(unknown driver) expected error")
	}
	if err.Error() == "" {
		t.Error("error message should be non-empty")
	}
}
