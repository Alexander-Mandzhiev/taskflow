package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	userCache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache"
	userRepository "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/user"
)

var _ repository.UserRepository = (*Adapter)(nil)

// Adapter — адаптер поверх UserReaderRepository, UserWriterRepository, UserCacheRepository.
type Adapter struct {
	reader userRepository.UserReaderRepository
	writer userRepository.UserWriterRepository
	cache  userCache.UserCacheRepository
}

// NewAdapter создаёт адаптер, объединяющий reader, writer и cache.
func NewAdapter(
	reader userRepository.UserReaderRepository,
	writer userRepository.UserWriterRepository,
	cache userCache.UserCacheRepository,
) *Adapter {
	return &Adapter{
		reader: reader,
		writer: writer,
		cache:  cache,
	}
}
