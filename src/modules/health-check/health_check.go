package health_check

import (
	"github.com/juanPabloCano/library-management-system/src/utils"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/health-check", run())
}

func run() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, utils.SuccessResponse{Message: "Server running ok"})
	}
}
