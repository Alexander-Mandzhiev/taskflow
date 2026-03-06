package http

import "time"

// Единый источник констант для пакета http (middleware: security, rate limit, body limit).

// Rate limit — общий максимум записей для IP и user rate limiter.
const RateLimiterMaxSize = 100_000

// Окно 5 мин: интервал очистки IP limiter и TTL записей user limiter.
const RateLimitWindow5Min = 5 * time.Minute

// Минимальное значение Retry-After (секунды) при 429 для IP и user rate limit.
const RateLimitRetryAfterSeconds = 1

// IP rate limiter
const (
	IPRateLimiterTTL              = 15 * time.Minute
	IPRateLimitDefaultPerSecond   = 30
	IPRateLimitSensitivePerSecond = 7
)

// User rate limiter
const (
	UserRateLimitWindow   = time.Minute
	UserRateLimitMax      = 100
	UserLimiterCleanupInt = 2 * time.Minute
)

// MaxRequestBodyBytes — лимит тела запроса для POST (1 MiB), защита от исчерпания памяти.
const MaxRequestBodyBytes = 1 << 20

// DefaultAccessTokenCookieName — имя cookie для JWT access-токена по умолчанию (fallback в JWTAuthMiddleware).
const DefaultAccessTokenCookieName = "access_token"
