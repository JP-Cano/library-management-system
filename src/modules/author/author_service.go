package author

import (
	"github.com/juanPabloCano/library-management-system/src/config"
	"github.com/juanPabloCano/library-management-system/src/database"
	"github.com/juanPabloCano/library-management-system/src/modules/author/entities"
)

type Repository interface {
	GetAll() ([]*entities.Author, error)
}

type Storage struct {
	db *database.Storage
}

func NewStorage(db *database.Storage) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAll() ([]*entities.Author, error) {
	rows, err := s.db.Conn.Query(config.BCGContext, `SELECT * FROM "library-management-system".public.authors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var authors []*entities.Author
	for rows.Next() {
		var author entities.Author
		if err := rows.Scan(&author.Id, &author.Name); err != nil {
			return nil, err
		}
		authors = append(authors, &author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}
