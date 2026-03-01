package password

import "golang.org/x/crypto/bcrypt"

const defaultCost = bcrypt.DefaultCost

// Hasher хеширует пароль для хранения. Compare проверяет пароль с хешем.
type Hasher interface {
	Hash(plain string) (string, error)
	Compare(hashed, plain string) error
}

// BcryptHasher реализация Hasher на bcrypt.
type BcryptHasher struct {
	Cost int
}

// NewBcryptHasher создаёт хешер с заданной сложностью (0 = default).
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost <= 0 {
		cost = defaultCost
	}
	return &BcryptHasher{Cost: cost}
}

// Hash возвращает bcrypt-хеш пароля.
func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.Cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare возвращает nil, если plain совпадает с hashed.
func (h *BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
