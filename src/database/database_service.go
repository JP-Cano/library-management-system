package database

import "github.com/jackc/pgx/v5"

type Storage struct {
	Conn *pgx.Conn
}

func NewStorage(conn *pgx.Conn) *Storage {
	return &Storage{
		Conn: conn,
	}
}
