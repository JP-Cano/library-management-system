package book

import (
	"github.com/google/uuid"
	"github.com/juanPabloCano/library-management-system/src/config"
	"github.com/juanPabloCano/library-management-system/src/database"
	"github.com/juanPabloCano/library-management-system/src/modules/book/entities"
)

type Repository interface {
	GetAllAvailable() ([]*entities.Book, error)
	FindByTitleOrAuthor(search string) ([]*entities.Book, error)
	UpdateStatus(id uuid.UUID, status bool) error
	GetById(id uuid.UUID) (*entities.Book, error)
}

type Storage struct {
	db *database.Storage
}

func NewStorage(db *database.Storage) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAllAvailable() ([]*entities.Book, error) {
	rows, err := s.db.Conn.Query(config.BCGContext, `
        SELECT books.id, books.title, books.available::boolean, authors.name AS author
        FROM "library-management-system".public.books
        JOIN "library-management-system".public.authors ON books.author_id = authors.id
        WHERE books.available = true
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entities.Book
	for rows.Next() {
		var book entities.Book
		if err := rows.Scan(&book.Id, &book.Title, &book.Available, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (s *Storage) FindByTitleOrAuthor(search string) ([]*entities.Book, error) {
	query := `
        SELECT books.id, books.title, books.available::boolean, authors.name AS author
        FROM "library-management-system".public.books
        JOIN "library-management-system".public.authors ON books.author_id = authors.id
        WHERE books.title ILIKE $1 OR authors.name ILIKE $2
    `
	rows, err := s.db.Conn.Query(config.BCGContext, query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entities.Book
	for rows.Next() {
		var book entities.Book
		if err := rows.Scan(&book.Id, &book.Title, &book.Available, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (s *Storage) UpdateStatus(id uuid.UUID, status bool) error {
	query := `
        UPDATE "library-management-system".public.books 
        SET available = $1 
        WHERE id = $2
    `
	_, err := s.db.Conn.Exec(config.BCGContext, query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetById(id uuid.UUID) (*entities.Book, error) {
	var b entities.Book
	var authorName string
	err := s.db.Conn.QueryRow(config.BCGContext, `
        SELECT books.id, books.title, books.available::boolean, authors.name AS author
        FROM "library-management-system".public.books
        JOIN "library-management-system".public.authors ON books.author_id = authors.id
        WHERE books.id = $1`, id).Scan(&b.Id, &b.Title, &b.Available, &authorName)
	if err != nil {
		return nil, err
	}
	b.Author = authorName
	return &b, nil
}
