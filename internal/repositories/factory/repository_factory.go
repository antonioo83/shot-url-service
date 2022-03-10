package factory

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/filestore"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
)

func GetRepository(config config.Config) interfaces.ShotURLRepository {
	if config.IsUseFileStore {
		return filestore.NewFileRepository(
			config.FileStoragePath,
		)
	} else {
		var buffer = make(map[string]models.ShortURL)
		return localcache.NewMemoryRepository(buffer)
	}
}
