package database

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/services"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"testing"
	"time"
)

const UserCode = 1
const ModelID = "1"

var pool *pgxpool.Pool
var rep interfaces.ShotURLRepository

func BenchmarkQueries(b *testing.B) {
	configFromFile, err := services.LoadConfigFile("config.json")
	if err != nil {
		fmt.Errorf("i can't load cofiguration file %w", err)
	}
	config := config.GetConfigSettings(configFromFile)
	context := context.Background()
	pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
	defer pool.Close()
	rep = NewShortURLRepository(context, pool)

	b.Run("saveUrl", func(b *testing.B) {
		urlModels := getRandomModels(UserCode, 1)
		urlModels[0].Code = ModelID
		for i := 0; i < b.N; i++ {
			rep.SaveURL(urlModels[0])
		}
	})

	b.Run("FindByCode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rep.FindByCode(ModelID)
		}
	})

	b.Run("SaveModels", func(b *testing.B) {
		urlModels := getRandomModels(UserCode, 10)
		for i := 0; i < b.N; i++ {
			rep.SaveModels(urlModels)
		}
	})

	b.Run("FindAllByUserCode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rep.FindAllByUserCode(UserCode)
		}
	})
}

func getRandomModels(userCode int, count int) []models.ShortURL {
	var urlModels []models.ShortURL
	for i := 0; i < count; i++ {
		model := models.ShortURL{}
		model.Code = strconv.Itoa(i) + "_" + time.Now().Format(time.RFC3339)
		model.UserCode = userCode
		model.CorrelationID = ""
		model.OriginalURL = "benchmark_" + time.Now().Format(time.RFC3339)
		model.ShortURL = "shotURL_" + strconv.Itoa(i) + "_" + time.Now().Format(time.StampNano)
		model.Active = false
		urlModels = append(urlModels, model)
	}

	return urlModels
}
