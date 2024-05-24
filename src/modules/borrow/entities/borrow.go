package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/juanPabloCano/library-management-system/src/modules/book/entities"
)

type Borrow struct {
	Id         uuid.UUID     `json:"id"`
	BookId     uuid.UUID     `json:"-"`
	UserId     uuid.UUID     `json:"-"`
	Book       entities.Book `json:"book"`
	User       string        `json:"user"`
	BorrowDate time.Time     `json:"borrow_date"`
	ReturnDate *time.Time    `json:"return_date"`
}

type BooksBorrowedToUser struct {
	Book       *entities.Book `json:"book"`
	BorrowDate time.Time      `json:"borrow_date"`
	ReturnDate *time.Time     `json:"return_date"`
}

type NewBorrow struct {
	BookId uuid.UUID `json:"bookId"`
	UserId uuid.UUID `json:"userId"`
}

type ReturnBorrow struct {
	BorrowId uuid.UUID `json:"borrowId"`
	BookId   uuid.UUID `json:"bookId"`
}
