package validation

import (
	"github.com/go-playground/validator/v10"
)

// Validator — общий экземпляр валидатора для проверки DTO в ручках (request body).
// Использование: validation.Validator.Struct(req)
var Validator = validator.New()
