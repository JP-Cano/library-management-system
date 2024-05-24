package book

import (
	"github.com/juanPabloCano/library-management-system/src/utils"
	"log"
	"net/http"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/books", s.GetAllAvailable())
	r.HandleFunc("/books/search", s.FindByTitleOrAuthor())
}

func (s *Service) GetAllAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		availableBooks, err := s.repository.GetAllAvailable()
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error getting all available books"})
			log.Println(err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, availableBooks)
	}
}

func (s *Service) FindByTitleOrAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("q")
		foundBooks, err := s.repository.FindByTitleOrAuthor(search)
		if len(foundBooks) == 0 {
			utils.WriteJSON(w, http.StatusOK, []any{})
			return
		}

		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error filtering books by title or author"})
			log.Println(err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, foundBooks)
	}
}
