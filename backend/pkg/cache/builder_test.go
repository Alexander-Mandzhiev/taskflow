package cache

import (
	"testing"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func TestBuildClient_EmptyAddr(t *testing.T) {
	client, err := BuildClient(&logger.NoopLogger{}, "", WithAddr(""))
	if err != nil {
		t.Fatalf("BuildClient(WithAddr(\"\")) err = %v", err)
	}
	if client != nil {
		t.Error("BuildClient(WithAddr(\"\")) client should be nil (cache disabled)")
	}
}

func TestBuildClient_NilLogger(t *testing.T) {
	// nil log заменяется на NoopLogger внутри BuildClient
	client, err := BuildClient(nil, "", WithAddr("localhost:6379"))
	if err != nil {
		t.Fatalf("BuildClient(nil log) err = %v", err)
	}
	if client == nil {
		t.Fatal("BuildClient(nil log) client should be non-nil")
	}
}
