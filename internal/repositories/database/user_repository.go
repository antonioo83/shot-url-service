package database

import (
	"context"
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	context    context.Context
	connection *pgxpool.Pool
}

func NewUserRepository(context context.Context, pool *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{context, pool}
}

//GetCount gets count of short url in the storage.
func (u userRepository) GetCount() (int, error) {
	var count int
	row := u.connection.QueryRow(u.context, "SELECT COUNT(*) FROM users")
	err := row.Scan(&count)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return count, nil
}

//Save saves a user in the storage.
func (u userRepository) Save(model models.User) error {
	_, err := u.connection.Query(u.context, "INSERT INTO users(code, uid)VALUES ($1, $2)", &model.Code, &model.UID)
	return err
}

//FindByCode finds a user in the storage by unique code.
func (u userRepository) FindByCode(code int) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(u.context, "SELECT code, uid FROM users WHERE code=$1", code).Scan(&model.Code, &model.UID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

//IsInDatabase check exists a user in the storage by unique code.
func (u userRepository) IsInDatabase(code int) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}

//GetLastModel gets a last user from the storage.
func (u userRepository) GetLastModel() (*models.User, error) {
	model := models.User{}
	err := u.connection.QueryRow(u.context, "SELECT code, uid FROM users ORDER BY code DESC").Scan(&model.Code, &model.UID)
	if errors.Is(err, pgx.ErrNoRows) {

		return &model, nil
	}

	return &model, err
}
