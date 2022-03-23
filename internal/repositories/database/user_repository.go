package database

import (
	"context"
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

func (u userRepository) Save(model models.User) error {
	_, err := u.connection.Query(u.context, "INSERT INTO users(code, uid)VALUES ($1, $2)", &model.CODE, &model.UID)
	return err
}

func (u userRepository) FindByCode(code int) (*models.User, error) {
	var model models.User
	err := u.connection.QueryRow(u.context, "SELECT code, uid FROM users WHERE code=$1", code).Scan(&model.CODE, &model.UID)
	if err == pgx.ErrNoRows {

		return nil, err
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u userRepository) IsInDatabase(code int) (bool, error) {
	model, err := u.FindByCode(code)

	return !(model == nil), err
}

func (u userRepository) GetLastModel() (*models.User, error) {
	model := models.User{}
	err := u.connection.QueryRow(u.context, "SELECT code, uid FROM users ORDER BY code DESC").Scan(&model.CODE, &model.UID)
	if err == pgx.ErrNoRows {

		return &model, err
	}

	return &model, nil
}
