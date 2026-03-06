package http

import (
	"net/http"
	"time"
)

// SetCookie устанавливает cookie с указанным именем, значением и TTL.
// httpOnly: true — cookie недоступна из JS (refresh); false — доступна (access token для фронта).
// domain — домен для cookie (например, ".example.com"). Пустой — только текущий домен.
func SetCookie(w http.ResponseWriter, name, value string, maxAgeSeconds int, isSecure bool, domain string, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: httpOnly,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
	}

	// Устанавливаем домен только если он указан
	if domain != "" {
		cookie.Domain = domain
	}

	if maxAgeSeconds > 0 {
		cookie.Expires = time.Now().Add(time.Duration(maxAgeSeconds) * time.Second)
	}

	http.SetCookie(w, cookie)
}

// DeleteCookie удаляет cookie с указанным именем
// Используется для очистки cookies при logout или истечении сессии
// domain - домен для cookie (например, ".classplanner.ru" для всех поддоменов). Если пустой, cookie доступна только для текущего домена
func DeleteCookie(w http.ResponseWriter, name string, isSecure bool, domain string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}

	// Устанавливаем домен только если он указан
	if domain != "" {
		cookie.Domain = domain
	}

	http.SetCookie(w, cookie)
}
