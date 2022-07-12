package interfaces

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DatabaseRepository interface {
	//Connect connects to the database.
	Connect(context context.Context) (*pgx.Conn, error)
	//Disconnect disconnects from the database.
	Disconnect(context context.Context, conn *pgx.Conn)
	//RunDump runs a sql dump by the filepath.
	RunDump(context context.Context, conn *pgxpool.Pool, filepath string) error
}
