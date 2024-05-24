package entities

import "github.com/google/uuid"

type Author struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
