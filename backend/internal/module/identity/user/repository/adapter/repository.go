package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	userCache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache"
	userRepository "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/user"
)

var _ repository.UserRepository = (*Repository)(nil)

// Repository — адаптер поверх UserReaderRepository, UserWriterRepository, UserCacheRepository.
type Repository struct {
	reader userRepository.UserReaderRepository
	writer userRepository.UserWriterRepository
	cache  userCache.UserCacheRepository
}

// NewRepository создаёт адаптер, объединяющий reader, writer и cache.
func NewRepository(
	reader userRepository.UserReaderRepository,
	writer userRepository.UserWriterRepository,
	cache userCache.UserCacheRepository,
) *Repository {
	return &Repository{
		reader: reader,
		writer: writer,
		cache:  cache,
	}
}
