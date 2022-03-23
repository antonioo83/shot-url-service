package database

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type databaseRepository struct {
	connString string
}

func NewDatabaseRepository(connString string) interfaces.DatabaseRepository {
	return &databaseRepository{connString}
}

func (r databaseRepository) Connect(context context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context, r.connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	//defer conn.Close(context)

	return conn, nil
}

func (r databaseRepository) Disconnect(context context.Context, conn *pgx.Conn) {
	defer conn.Close(context)
}

func (r databaseRepository) RunDump(context context.Context, conn *pgxpool.Pool, filepath string) error {
	sqlDump, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context, string(sqlDump))

	return err
}
