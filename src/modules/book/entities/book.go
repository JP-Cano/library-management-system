package entities

import "github.com/google/uuid"

type Book struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Available bool      `json:"available"`
	Author    string    `json:"author"`
	AuthorId  string    `json:"-"`
}
