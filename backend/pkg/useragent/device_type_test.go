package useragent

import (
	"testing"
)

func TestDeviceTypeFromUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		want     string
	}{
		{"empty", "", DeviceTypeUnknown},
		{"whitespace only", "   \t  ", DeviceTypeUnknown},
		{"tablet", "tablet", DeviceTypeTablet},
		{"Tablet uppercase", "Tablet", DeviceTypeTablet},
		{"ipad", "ipad", DeviceTypeTablet},
		{"iPad in string", "Mozilla/5.0 iPad", DeviceTypeTablet},
		{"mobile", "mobile", DeviceTypeMobile},
		{"android", "android", DeviceTypeMobile},
		{"Android in string", "Mozilla/5.0 Android", DeviceTypeMobile},
		{"iphone", "iphone", DeviceTypeMobile},
		{"iPhone in string", "Mozilla/5.0 iPhone", DeviceTypeMobile},
		{"desktop default", "Mozilla/5.0 Windows NT", DeviceTypeDesktop},
		{"desktop Chrome", "Chrome/120.0 Windows", DeviceTypeDesktop},
		{"trimmed and lower", "  ANDROID  ", DeviceTypeMobile},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeviceTypeFromUserAgent(tt.ua)
			if got != tt.want {
				t.Errorf("DeviceTypeFromUserAgent(%q) = %q, want %q", tt.ua, got, tt.want)
			}
		})
	}
}

func TestDeviceType_constants(t *testing.T) {
	if DeviceTypeDesktop != "desktop" {
		t.Errorf("DeviceTypeDesktop = %q, want desktop", DeviceTypeDesktop)
	}
	if DeviceTypeMobile != "mobile" {
		t.Errorf("DeviceTypeMobile = %q, want mobile", DeviceTypeMobile)
	}
	if DeviceTypeTablet != "tablet" {
		t.Errorf("DeviceTypeTablet = %q, want tablet", DeviceTypeTablet)
	}
	if DeviceTypeUnknown != "unknown" {
		t.Errorf("DeviceTypeUnknown = %q, want unknown", DeviceTypeUnknown)
	}
}
