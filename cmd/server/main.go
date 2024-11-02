package main

import (
	"log/slog"
	"main/internal/config"
	ApplicationHandler "main/internal/handlers"
	"main/internal/lib/liblogger"
	"main/internal/middleware/midlogger"
	ApplicationRepository "main/internal/repositories"
	ApplicationService "main/internal/services"
	"main/internal/storage/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const LocalFilePath = "D:/go_dev/Olimpiad_portal/Application_Service/config-yaml/local.yaml"

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
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)

	// init application service and handler
	applicationService := ApplicationService.NewApplicationService(storage, &ApplicationRepository.ApplicationRepository{})
	applicationHandler := ApplicationHandler.NewApplicationHandler(applicationService, log)

	// init applications route
	router.Get("/applications", applicationHandler.GetAllApplications)
	router.Get("/applications/{id}", applicationHandler.GetApplicationByID)
	router.Post("/applications", applicationHandler.CreateApplication)
	router.Put("/applications/{id}", applicationHandler.UpdateApplicationStatus)
	router.Delete("/applications/{id}", applicationHandler.DeleteApplication)

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
