package model

import (
	"time"

	"github.com/google/uuid"
)

// Session — данные сессии в кеше (ключ session:{session_id}).
// Метадата используется для списка активных сессий и отображения устройства.
type Session struct {
	UserID     uuid.UUID
	CreatedAt  time.Time
	DeviceType string // useragent.DeviceTypeDesktop | DeviceTypeMobile | DeviceTypeTablet | DeviceTypeUnknown
	UserAgent  string // сырой User-Agent при логине (опционально)
	IP         string // IP при логине (опционально)
}
