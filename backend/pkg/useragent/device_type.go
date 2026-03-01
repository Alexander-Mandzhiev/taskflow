package useragent

import "strings"

// Типы устройств для сессии (список сессий, отображение "с какого устройства зашли").
const (
	DeviceTypeDesktop = "desktop"
	DeviceTypeMobile  = "mobile"
	DeviceTypeTablet  = "tablet"
	DeviceTypeUnknown = "unknown"
)

// DeviceTypeFromUserAgent возвращает константу типа устройства по строке User-Agent.
// Пустой userAgent → DeviceTypeUnknown. Используется при создании сессии (Login).
func DeviceTypeFromUserAgent(ua string) string {
	lower := strings.ToLower(strings.TrimSpace(ua))
	switch {
	case lower == "":
		return DeviceTypeUnknown
	case strings.Contains(lower, "tablet") || strings.Contains(lower, "ipad"):
		return DeviceTypeTablet
	case strings.Contains(lower, "mobile") || strings.Contains(lower, "android") || strings.Contains(lower, "iphone"):
		return DeviceTypeMobile
	default:
		return DeviceTypeDesktop
	}
}
