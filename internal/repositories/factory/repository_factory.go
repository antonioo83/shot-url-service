package factory

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/filestore"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
	"log"
)

func GetRepository(config config.Config) interfaces.ShotURLRepository {
	if config.IsUseFileStore {
		consumer, err := filestore.GetConsumer(config.FileStoragePath)
		if err != nil {
			log.Fatalln("can't create consumer: ", err)
		}

		producer, err := filestore.GetProducer(config.FileStoragePath)
		if err != nil {
			log.Fatalln("can't create producer: ", err)
		}

		return filestore.NewFileRepository(*consumer, *producer)
	} else {
		var buffer = make(map[string]models.ShortURL)
		return localcache.NewMemoryRepository(buffer)
	}
}
