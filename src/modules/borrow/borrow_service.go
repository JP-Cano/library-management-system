package borrow

import (
	"errors"
	"github.com/google/uuid"
	"github.com/juanPabloCano/library-management-system/src/config"
	"github.com/juanPabloCano/library-management-system/src/database"
	"github.com/juanPabloCano/library-management-system/src/modules/book"
	"github.com/juanPabloCano/library-management-system/src/modules/borrow/entities"
	"github.com/juanPabloCano/library-management-system/src/modules/user"
	"github.com/juanPabloCano/library-management-system/src/utils"
)

type Repository interface {
	BooksBorrowedToUsers(userId uuid.UUID) ([]*entities.BooksBorrowedToUser, error)
	Register(bookId, userId uuid.UUID) (*entities.Borrow, error)
	Return(borrowId, bookId uuid.UUID) (*utils.SuccessResponse, error)
}

type Storage struct {
	db             *database.Storage
	bookRepository book.Repository
	userRepository user.Repository
}

const (
	available   = true
	unavailable = false
)

func NewStorage(db *database.Storage, bookRepository book.Repository, userRepository user.Repository) *Storage {
	return &Storage{db: db, bookRepository: bookRepository, userRepository: userRepository}
}

func (s *Storage) BooksBorrowedToUsers(userId uuid.UUID) ([]*entities.BooksBorrowedToUser, error) {
	borrows, err := s.getBorrowsByUserID(userId)
	if err != nil {
		return nil, err
	}

	var booksBorrowedToUsers []*entities.BooksBorrowedToUser
	for _, borrow := range borrows {
		foundBook, err := s.bookRepository.GetById(borrow.BookId)
		if err != nil {
			return nil, err
		}
		borrowedWithBook := &entities.BooksBorrowedToUser{
			Book:       foundBook,
			BorrowDate: borrow.BorrowDate,
			ReturnDate: borrow.ReturnDate,
		}
		booksBorrowedToUsers = append(booksBorrowedToUsers, borrowedWithBook)
	}

	return booksBorrowedToUsers, nil
}

func (s *Storage) getBorrowsByUserID(userId uuid.UUID) ([]*entities.Borrow, error) {
	var borrows []*entities.Borrow
	query := `SELECT id, book_id, user_id, borrow_date, return_date FROM "library-management-system".public.borrows WHERE user_id = $1`
	rows, err := s.db.Conn.Query(config.BCGContext, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b entities.Borrow
		if err := rows.Scan(&b.Id, &b.BookId, &b.UserId, &b.BorrowDate, &b.ReturnDate); err != nil {
			return nil, err
		}
		borrows = append(borrows, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return borrows, nil
}

func (s *Storage) Register(bookId, userId uuid.UUID) (*entities.Borrow, error) {
	query := `
        INSERT INTO "library-management-system".public.borrows (book_id, user_id)
        VALUES ($1, $2)
        RETURNING id, book_id, user_id, borrow_date, return_date
    `
	selectedBook, selectedBookError := s.bookRepository.GetById(bookId)
	if !selectedBook.Available {
		return nil, errors.New("book is currently unavailable")
	}

	var borrow entities.Borrow
	err := s.db.Conn.QueryRow(config.BCGContext, query, bookId, userId).Scan(&borrow.Id, &borrow.BookId, &borrow.UserId, &borrow.BorrowDate, &borrow.ReturnDate)

	if err != nil {
		return nil, err
	}

	bookStatusError := s.bookRepository.UpdateStatus(selectedBook.Id, unavailable)

	if selectedBookError != nil {
		return nil, selectedBookError
	}

	if bookStatusError != nil {
		return nil, selectedBookError
	}

	selectedUser, selectedBookError := s.userRepository.GetById(userId)
	if selectedBookError != nil {
		return nil, selectedBookError
	}

	borrow.Book = *selectedBook
	borrow.User = selectedUser.Name

	return &borrow, nil
}

func (s *Storage) Return(borrowId, bookId uuid.UUID) (*utils.SuccessResponse, error) {
	query := `UPDATE "library-management-system".public.borrows SET return_date = NOW() WHERE id = $1`

	_, err := s.db.Conn.Exec(config.BCGContext, query, borrowId)
	if err != nil {
		return nil, err
	}
	err = s.bookRepository.UpdateStatus(bookId, available)
	if err != nil {
		return nil, err
	}

	return &utils.SuccessResponse{Message: "Book returned successfully"}, nil
}
