package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type PsqlStorage struct {
	Conn *pgx.Conn
}

var BCGContext = context.Background()

func NewPsqlStorage(url string) *PsqlStorage {
	conn, err := pgx.Connect(BCGContext, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(BCGContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database connected successfully")
	return &PsqlStorage{Conn: conn}
}

func (p *PsqlStorage) Init() (*pgx.Conn, error) {
	if err := p.createUUIDExtension(); err != nil {
		return nil, err
	}

	if err := p.createBooksTable(); err != nil {
		return nil, err
	}

	if err := p.createBorrowsTable(); err != nil {
		return nil, err
	}

	if err := p.createAuthorsTable(); err != nil {
		return nil, err
	}

	if err := p.createUsersTable(); err != nil {
		return nil, err
	}

	return p.Conn, nil
}

func (p *PsqlStorage) Close() {
	err := p.Conn.Close(BCGContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error closing database connection: %v\n", err)
	}
}

func (p *PsqlStorage) createUUIDExtension() error {
	_, err := p.Conn.Exec(BCGContext, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	return err
}

func (p *PsqlStorage) createBooksTable() error {
	_, err := p.Conn.Exec(BCGContext, `
		CREATE TABLE IF NOT EXISTS books (
			id CHAR(36) DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author_id CHAR(36) REFERENCES authors(id),
			available BOOLEAN
		);`)
	return err
}

func (p *PsqlStorage) createAuthorsTable() error {
	_, err := p.Conn.Exec(BCGContext, `
		CREATE TABLE IF NOT EXISTS authors (
			id CHAR(36) DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		);`)
	return err
}

func (p *PsqlStorage) createUsersTable() error {
	_, err := p.Conn.Exec(BCGContext, `
		CREATE TABLE IF NOT EXISTS users (
			id CHAR(36) DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL
		);`)
	return err
}

func (p *PsqlStorage) createBorrowsTable() error {
	_, err := p.Conn.Exec(BCGContext, `
		CREATE TABLE IF NOT EXISTS borrows (
			id CHAR(36) DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
			book_id CHAR(36) REFERENCES books(id) NOT NULL,
			user_id CHAR(36) REFERENCES users(id) NOT NULL,
			borrow_date DATE DEFAULT NOW() NOT NULL,
			return_date DATE
		);`)
	return err
}
