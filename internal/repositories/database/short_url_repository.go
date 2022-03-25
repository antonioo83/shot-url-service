package database

import (
	"context"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type shortUrlRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewShortUrlRepository(context context.Context, pool *pgxpool.Pool) interfaces.ShotURLRepository {
	return &shortUrlRepository{context, pool}
}

func (s shortUrlRepository) SaveURL(model models.ShortURL) error {
	_, err := s.connection.Exec(
		s.context,
		"INSERT INTO short_url(correlation_id, user_code, code, original_url, short_url)VALUES ($1, $2, $3, $4, $5)",
		model.CorrelationId, model.UserCode, model.Code, model.OriginalURL, model.ShortURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s shortUrlRepository) SaveModels(models map[int]models.ShortURL) error {
	tx, err := s.connection.Begin(s.context)
	if err != nil {
		panic(err)
	}

	b := &pgx.Batch{}
	for _, model := range models {
		b.Queue(
			"INSERT INTO short_url(correlation_id, user_code, code, original_url, short_url)VALUES ($1, $2, $3, $4, $5)",
			model.CorrelationId, model.UserCode, model.Code, model.OriginalURL, model.ShortURL,
		)
	}
	tx.SendBatch(s.context, b)

	return tx.Commit(s.context)
}

func (s shortUrlRepository) FindByCode(code string) (*models.ShortURL, error) {
	var model models.ShortURL
	row := s.connection.QueryRow(s.context, "SELECT correlation_id, user_code, code, original_url, short_url FROM short_url WHERE code=$1", code)
	err := row.Scan(&model.CorrelationId, &model.UserCode, &model.Code, &model.OriginalURL, &model.ShortURL)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (s shortUrlRepository) FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error) {
	var model = models.ShortURL{}
	models := make(map[string]models.ShortURL)
	rows, _ := s.connection.Query(s.context, "SELECT correlation_id, user_code, code, original_url, short_url FROM short_url WHERE user_code=$1", userCode)
	for rows.Next() {
		rows.Scan(&model.CorrelationId, &model.UserCode, &model.Code, &model.OriginalURL, &model.ShortURL)
		models[model.Code] = model
	}

	return &models, nil
}

func (s shortUrlRepository) IsInDatabase(code string) (bool, error) {
	model, err := s.FindByCode(code)

	return !(model == nil), err
}
