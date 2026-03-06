package txmanager

import (
	"testing"
)

func TestTruncateStr(t *testing.T) {
	tests := []struct {
		s    string
		max  int
		want string
	}{
		{"short", 10, "short"},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcdef", 3, "abc..."},
		{"long string here", 8, "long str..."},
	}
	for _, tt := range tests {
		got := truncateStr(tt.s, tt.max)
		if got != tt.want {
			t.Errorf("truncateStr(%q, %d) = %q, want %q", tt.s, tt.max, got, tt.want)
		}
	}
}

func TestFormatPanic(t *testing.T) {
	if got := formatPanic(nil); got != "" {
		t.Errorf("formatPanic(nil) = %q, want \"\"", got)
	}
	got := formatPanic("oops")
	if got == "" {
		t.Error("formatPanic(string) should be non-empty")
	}
	if len(got) > maxAttrLen+3 {
		t.Errorf("formatPanic result length %d > maxAttrLen+3", len(got))
	}
}
