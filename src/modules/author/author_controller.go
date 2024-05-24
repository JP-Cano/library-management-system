package author

import (
	"github.com/juanPabloCano/library-management-system/src/utils"
	"net/http"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/authors", s.getAll())
}

func (s *Service) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors, err := s.repository.GetAll()
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Error getting all authors"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, authors)
	}
}
