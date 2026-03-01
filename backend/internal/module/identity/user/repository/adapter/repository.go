package adapter

import (
	"mkk/internal/module/identity/user/repository"
)

var _ repository.UserRepository = (*Repository)(nil)

// Repository — адаптер поверх UserReaderRepository, UserWriterRepository, UserCacheRepository.
type Repository struct {
	reader repository.UserReaderRepository
	writer repository.UserWriterRepository
	cache  repository.UserCacheRepository
}

// NewRepository создаёт адаптер, объединяющий reader, writer и cache.
func NewRepository(
	reader repository.UserReaderRepository,
	writer repository.UserWriterRepository,
	cache repository.UserCacheRepository,
) *Repository {
	return &Repository{
		reader: reader,
		writer: writer,
		cache:  cache,
	}
}
