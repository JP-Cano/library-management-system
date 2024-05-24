package api

import (
	"log"
	"net/http"

	"github.com/juanPabloCano/library-management-system/src/database"
	"github.com/juanPabloCano/library-management-system/src/modules/author"
	"github.com/juanPabloCano/library-management-system/src/modules/book"
	"github.com/juanPabloCano/library-management-system/src/modules/borrow"
	"github.com/juanPabloCano/library-management-system/src/modules/health-check"
	"github.com/juanPabloCano/library-management-system/src/modules/user"
	"github.com/juanPabloCano/library-management-system/src/utils/middlewares"
)

type Server struct {
	Addr string
	DB   *database.Storage
}

func New(addr string, db *database.Storage) *Server {
	return &Server{
		Addr: addr,
		DB:   db,
	}
}

func (s *Server) Serve() error {
	router := http.NewServeMux()
	registerDependencies(s.DB, router)
	server := &http.Server{
		Addr:    s.Addr,
		Handler: middlewares.RequestLoggerMiddleware(router),
	}

	log.Printf("Server running on port %s", s.Addr)
	return server.ListenAndServe()
}

func registerDependencies(db *database.Storage, router *http.ServeMux) {
	// Health-check
	health_check.RegisterRoutes(router)

	// Authors
	authorRepo := author.NewStorage(db)
	authorService := author.NewService(authorRepo)
	authorService.RegisterRoutes(router)

	// Books
	bookRepo := book.NewStorage(db)
	bookService := book.NewService(bookRepo)
	bookService.RegisterRoutes(router)

	// Users
	userRepo := user.NewStorage(db)

	// Borrows
	borrowRepo := borrow.NewStorage(db, bookRepo, userRepo)
	borrowService := borrow.NewService(borrowRepo)
	borrowService.RegisterRoutes(router)
}
