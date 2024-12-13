package main

import (
	"OlimpiadPortal/ApplicationService/internal/config"
	ApplicationHandler "OlimpiadPortal/ApplicationService/internal/handlers"
	"OlimpiadPortal/ApplicationService/internal/lib/liblogger"
	"OlimpiadPortal/ApplicationService/internal/middleware/auth"
	"OlimpiadPortal/ApplicationService/internal/middleware/midlogger"
	ApplicationRepository "OlimpiadPortal/ApplicationService/internal/repositories"
	ApplicationService "OlimpiadPortal/ApplicationService/internal/services"
	"OlimpiadPortal/ApplicationService/internal/storage/postgresql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const LocalFilePath = "config-yaml/local.yaml"

func main() {
	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger
	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start application-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := postgresql.NewPostgreSQL(cfg.GetDataSourceName())
	if err != nil {
		log.Error("failed to init storage", liblogger.Err(err))
	} else {
		log.Info("storage is enabled")
	}

	// init chi router
	router := chi.NewRouter()
	// init cors
	corsOptions := cors.Options{
		AllowedOrigins:   []string{cfg.ReactVision}, // React URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // В секундах
	}
	router.Use(cors.Handler(corsOptions))
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)

	// add Authentication with JWT token
	router.Use(func(next http.Handler) http.Handler {
		return auth.AuthenticateMiddleware(next, cfg.Key)
	})

	// init application service and handler
	applicationService := ApplicationService.NewApplicationService(storage, &ApplicationRepository.ApplicationRepository{})
	applicationHandler := ApplicationHandler.NewApplicationHandler(applicationService, log)

	router.Get("/applications/{id}", applicationHandler.GetApplicationByID) //wtf

	// init applications route
	router.With(auth.RoleBasedAccess("4")).Group(func(r chi.Router) {
		router.Get("/applications", applicationHandler.GetAllApplications)
	})

	router.With(auth.RoleBasedAccess("2", "4")).Group(func(r chi.Router) {
		router.Get("/applications/user/{userID}", applicationHandler.GetApplicationsByUserID)
	})

	router.With(auth.RoleBasedAccess("3")).Group(func(r chi.Router) {
		router.Get("/applications/event/{eventID}", applicationHandler.GetApplicationsByEventID)
	})

	router.With(auth.RoleBasedAccess("2", "4")).Group(func(r chi.Router) {
		router.Post("/applications", applicationHandler.CreateApplication)
	})

	router.With(auth.RoleBasedAccess("3", "4")).Group(func(r chi.Router) {
		router.Put("/applications/{id}", applicationHandler.UpdateApplicationStatus)
		router.Delete("/applications/{id}", applicationHandler.DeleteApplication)
	})

	// init server
	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}
	// start server
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", liblogger.Err(err))
	}

	log.Error("server stopped")
}
