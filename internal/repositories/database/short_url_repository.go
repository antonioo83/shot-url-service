package database

import (
	"context"
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type shortURLRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewShortURLRepository(context context.Context, pool *pgxpool.Pool) interfaces.ShotURLRepository {
	return &shortURLRepository{context, pool}
}

//GetCount gets count of short url in the storage.
func (s shortURLRepository) GetCount() (int, error) {
	var count int
	row := s.connection.QueryRow(s.context, "SELECT COUNT(*) FROM short_url WHERE active=true")
	err := row.Scan(&count)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return count, nil
}

//SaveURL saves an entity in the storage.
func (s shortURLRepository) SaveURL(model models.ShortURL) error {
	_, err := s.connection.Exec(
		s.context,
		"INSERT INTO short_url(correlation_id, user_code, code, original_url, short_url, active)VALUES ($1, $2, $3, $4, $5, $6)",
		model.CorrelationID, model.UserCode, model.Code, model.OriginalURL, model.ShortURL, model.Active,
	)
	if err != nil {
		return err
	}

	return nil
}

//SaveModels saves batch of entities in the storage.
func (s shortURLRepository) SaveModels(models []models.ShortURL) error {
	tx, err := s.connection.BeginTx(s.context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(s.context)
		} else {
			tx.Commit(s.context)
		}
	}()

	for _, model := range models {
		_, err = tx.Exec(
			s.context,
			"INSERT INTO short_url(correlation_id, user_code, code, original_url, short_url, active)VALUES ($1, $2, $3, $4, $5, $6)",
			model.CorrelationID, model.UserCode, model.Code, model.OriginalURL, model.ShortURL, model.Active,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

//FindByCode finds an entity in the storage by unique code.
func (s shortURLRepository) FindByCode(code string) (*models.ShortURL, error) {
	var model models.ShortURL
	row := s.connection.QueryRow(s.context, "SELECT correlation_id, user_code, code, original_url, short_url, active FROM short_url WHERE code=$1", code)
	err := row.Scan(&model.CorrelationID, &model.UserCode, &model.Code, &model.OriginalURL, &model.ShortURL, &model.Active)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

//FindAllByUserCode finds entities in the storage by unique codes.
func (s shortURLRepository) FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error) {
	var model = models.ShortURL{}
	models := make(map[string]models.ShortURL)
	rows, err := s.connection.Query(s.context, "SELECT correlation_id, user_code, code, original_url, short_url, active FROM short_url WHERE user_code=$1", userCode)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&model.CorrelationID, &model.UserCode, &model.Code, &model.OriginalURL, &model.ShortURL, &model.Active)
		if err != nil {
			return nil, err
		}
		models[model.Code] = model
	}

	return &models, nil
}

//IsInDatabase check exists an entity in the storage by unique code.
func (s shortURLRepository) IsInDatabase(code string) (bool, error) {
	model, err := s.FindByCode(code)

	return !(model == nil), err
}

//Delete deletes entities of a user from the storage by user code.
func (s shortURLRepository) Delete(userCode int, codes []string) error {
	batch := &pgx.Batch{}
	for _, code := range codes {
		batch.Queue("UPDATE short_url SET active=false WHERE user_code=$1 AND code=$2", userCode, code)
	}
	br := s.connection.SendBatch(context.Background(), batch)

	_, err := br.Exec()

	return err
}
