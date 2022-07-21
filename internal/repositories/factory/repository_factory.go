package factory

import (
	"context"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/database"
	"github.com/antonioo83/shot-url-service/internal/repositories/filestore"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetRepository(context context.Context, pool *pgxpool.Pool, config config.Config) interfaces.ShotURLRepository {
	if config.IsUseDatabase {
		return database.NewShortURLRepository(context, pool)
	} else if config.IsUseFileStore {
		return filestore.NewFileRepository(config.FileStoragePath)
	} else {
		var buffer = make(map[string]models.ShortURL)
		return localcache.NewMemoryRepository(buffer)
	}
}

func GetUserRepository(context context.Context, pool *pgxpool.Pool, config config.Config) interfaces.UserRepository {
	if config.IsUseDatabase {
		return database.NewUserRepository(context, pool)
	} else if config.IsUseFileStore {
		return filestore.NewUserRepository(config.UserFileStoragePath)
	} else {
		var userBuffer = make(map[int]models.User)
		return localcache.NewMemoryUserRepository(userBuffer)
	}
}

func GetDatabaseRepository(config config.Config) interfaces.DatabaseRepository {
	return database.NewDatabaseRepository(config.DatabaseDsn)
}
