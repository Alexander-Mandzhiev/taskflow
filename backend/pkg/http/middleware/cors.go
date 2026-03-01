package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSMiddleware создаёт middleware для обработки CORS запросов.
// Автоматически обрабатывает preflight запросы и устанавливает необходимые заголовки.
// Поддерживает wildcard origins (*) и точное соответствие доменов.
func CORSMiddleware(
	allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders []string,
	allowCredentials bool,
	maxAge int,
) func(http.Handler) http.Handler {
	allowMethods := strings.Join(allowedMethods, ", ")
	allowHeaders := strings.Join(allowedHeaders, ", ")
	exposeHeaders := strings.Join(exposedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && (len(allowedOrigins) == 0 || matchesOrigin(origin, allowedOrigins)) {
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Origin", origin)

				if allowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if allowMethods != "" {
					w.Header().Set("Access-Control-Allow-Methods", allowMethods)
				}

				if allowHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				}

				if exposeHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
				}

				if maxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(maxAge))
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// matchesOrigin проверяет, разрешён ли данный origin в списке разрешённых.
// Поддерживает wildcard (*) для разрешения всех доменов.
func matchesOrigin(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}

	return false
}
