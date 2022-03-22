package interfaces

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DatabaseRepository interface {
	Connect(context context.Context) (*pgx.Conn, error)
	Disconnect(context context.Context, conn *pgx.Conn)
}
