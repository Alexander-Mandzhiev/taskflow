package util

import (
	"testing"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"valid", "42", 42, false},
		{"zero", "0", 0, false},
		{"negative", "-10", -10, false},
		{"invalid", "abc", 0, true},
		{"float_string", "3.14", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseInt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseInt_emptyStringError(t *testing.T) {
	_, err := ParseInt("")
	if err == nil {
		t.Fatal("ParseInt(\"\") expected error")
	}
	if err.Error() != "empty string" {
		t.Errorf("ParseInt(\"\") error = %q, want \"empty string\"", err.Error())
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"valid", "9223372036854775807", int64(9223372036854775807), false},
		{"zero", "0", 0, false},
		{"negative", "-1", -1, false},
		{"invalid", "xyz", 0, true},
		{"overflow_positive", "99999999999999999999", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt64(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseInt64(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{"empty", "", false, true},
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"invalid", "yes", false, true},
		{"invalid_num", "2", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBool(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseBool(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
