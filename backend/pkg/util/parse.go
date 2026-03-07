// Package util содержит общие утилиты (парсинг строк и т.д.).
package util

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

// ParseInt парсит строку в int.
// Возвращает ошибку, если строка пустая или не может быть преобразована в int.
func ParseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int: %w", err)
	}
	return result, nil
}

// ParseInt64 парсит строку в int64 (для UnixNano timestamps).
// Возвращает ошибку, если строка пустая или не может быть преобразована в int64.
func ParseInt64(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int64: %w", err)
	}
	return result, nil
}

// ParseBool парсит строку в bool.
// Возвращает ошибку, если строка пустая или не может быть преобразована в bool.
func ParseBool(s string) (bool, error) {
	if s == "" {
		return false, fmt.Errorf("empty string")
	}
	result, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("failed to parse bool: %w", err)
	}
	return result, nil
}

// ParseUUID парсит строку в uuid.UUID (для id в URL и т.д.).
// Возвращает ошибку, если строка пустая или не является валидным UUID.
func ParseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, fmt.Errorf("empty string")
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse uuid: %w", err)
	}
	return parsed, nil
}
