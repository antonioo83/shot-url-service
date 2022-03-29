package interfaces

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DatabaseRepository interface {
	Connect(context context.Context) (*pgx.Conn, error)
	Disconnect(context context.Context, conn *pgx.Conn)
	RunDump(context context.Context, conn *pgxpool.Pool, filepath string) error
}
