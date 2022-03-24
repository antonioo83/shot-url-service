package database

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	sqlDump := "CREATE TABLE IF NOT EXISTS short_url (id serial NOT NULL PRIMARY KEY,user_code integer NOT NULL,code character varying(50) NOT NULL,original_url character varying(500) NOT NULL,short_url character varying(500) NOT NULL);CREATE TABLE IF NOT EXISTS users (id serial NOT NULL PRIMARY KEY,code integer NOT NULL,uid character varying(500) NOT NULL);"
	_, err := conn.Exec(context, sqlDump)

	return err
}
