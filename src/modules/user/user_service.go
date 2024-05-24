package user

import (
	"github.com/google/uuid"
	"github.com/juanPabloCano/library-management-system/src/config"
	"github.com/juanPabloCano/library-management-system/src/database"
	"github.com/juanPabloCano/library-management-system/src/modules/user/entities"
)

type Repository interface {
	GetById(id uuid.UUID) (*entities.User, error)
}

type Storage struct {
	db *database.Storage
}

func NewStorage(db *database.Storage) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetById(id uuid.UUID) (*entities.User, error) {
	var u entities.User
	err := s.db.Conn.QueryRow(config.BCGContext, `SELECT * FROM "library-management-system".public.users WHERE users.id = $1`, id).Scan(&u.Id, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
