package task_v1

import (
	"strconv"
)

// parseIntPositive парсит положительное целое из строки. При ошибке или n <= 0 возвращает (0, error).
func parseIntPositive(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, strconv.ErrSyntax
	}
	return n, nil
}
