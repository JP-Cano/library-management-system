package borrow

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/juanPabloCano/library-management-system/src/modules/borrow/entities"
	"github.com/juanPabloCano/library-management-system/src/utils"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /borrows/search", s.GetBooksBorrowedToUser())
	r.HandleFunc("POST /borrows/register", s.Register())
	r.HandleFunc("POST /borrows/return", s.Return())
}

func (s *Service) GetBooksBorrowedToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("u")
		booksLoanedToUser, err := s.repository.BooksBorrowedToUsers(uuid.MustParse(userId))
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error getting all borrowed books by user"})
			log.Println(err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, booksLoanedToUser)
	}
}

func (s *Service) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
		}

		defer r.Body.Close()
		var payload *entities.NewBorrow
		err = json.Unmarshal(body, &payload)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error parsing body"})
		}

		borrow, err := s.repository.Register(payload.BookId, payload.UserId)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
			log.Println(err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, borrow)
	}
}

func (s *Service) Return() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
		}

		defer r.Body.Close()
		var payload *entities.ReturnBorrow
		err = json.Unmarshal(body, &payload)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error parsing body"})
		}
		response, err := s.repository.Return(payload.BorrowId, payload.BookId)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
			log.Println(err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, response)
	}
}
