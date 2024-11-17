package server

import (
	"context"
	"fmt"
	authInterface "image-processing-service/src/internal/auth/interfaces"
	imageInterface "image-processing-service/src/internal/images/interfaces"
	userInterface "image-processing-service/src/internal/users/interfaces"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	port        int
	server      *http.Server
	authServer  *authInterface.AuthServer
	userServer  *userInterface.UserServer
	imageServer *imageInterface.ImageServer
}

func NewService(port int, authService *authInterface.AuthServer, usersCfg *userInterface.UserServer, imagesCfg *imageInterface.ImageServer) *Service {
	service := &Service{
		port:        port,
		authServer:  authService,
		userServer:  usersCfg,
		imageServer: imagesCfg,
	}
	service.setup()
	return service
}

func (s *Service) setup() {
	mux := http.NewServeMux()
	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           mux,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("/health", health)

	mux.HandleFunc("POST /users", s.userServer.Register)
	mux.HandleFunc("GET /users", s.authServer.UserMiddleware(s.userServer.Info))
	mux.HandleFunc("PUT /users", s.authServer.UserMiddleware(s.userServer.Update))
	mux.HandleFunc("DELETE /users", s.authServer.UserMiddleware(s.userServer.Delete))

	mux.HandleFunc("POST /login", s.authServer.Login)
	mux.HandleFunc("POST /refresh", s.authServer.Refresh)
	mux.HandleFunc("DELETE /logout", s.authServer.UserMiddleware(s.authServer.Logout))

	mux.HandleFunc("POST /images", s.authServer.UserMiddleware(s.imageServer.Upload))
	mux.HandleFunc("GET /images/file", s.authServer.UserMiddleware(s.imageServer.Download))
	mux.HandleFunc("GET /images/info", s.authServer.UserMiddleware(s.imageServer.Info))
	mux.HandleFunc("GET /images/list", s.authServer.UserMiddleware(s.imageServer.List))
	mux.HandleFunc("PUT /images", s.authServer.UserMiddleware(s.imageServer.Transform))
	mux.HandleFunc("DELETE /images", s.authServer.UserMiddleware(s.imageServer.Delete))

	mux.HandleFunc("GET /admin/users", s.authServer.AdminMiddleware(s.userServer.AdminListAllUsers))
	mux.HandleFunc("PATCH /admin/users", s.authServer.AdminMiddleware(s.userServer.AdminUpdateRole))
	mux.HandleFunc("DELETE /admin/users", s.authServer.AdminMiddleware(s.userServer.AdminDeleteUser))
}

func (s *Service) Start() {
	slog.Info("Starting server", "port", s.port)
	err := s.server.ListenAndServe()
	if err != nil {
		slog.Error("Init error: error starting server", "error", err)
		return
	}
}

func (s *Service) Stop() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		slog.Error("Shutdown error: error shutting down server", "error", err)
	}
	slog.Info("Server stopped")
}
