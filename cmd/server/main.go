package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/subject_handler"
	"main/internal/lib/liblogger"
	"main/internal/middleware/midlogger"
	"main/internal/repositories/subject_repository"
	"main/internal/services/subject_service"
	"main/internal/storage/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const LocalFilePath = "config-yaml/local.yaml"

func main() {
	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger
	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start event-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := postgresql.NewPosgreSQL(cfg.GetDataSourceName())
	if err != nil {
		log.Error("faile to init storage", liblogger.Err(err))
	} else {
		log.Info("storage are enabled")
	}

	// init chi router
	router := chi.NewRouter()
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)
	// init subject service and handler
	subject_service := subject_service.NewSubjectService(storage, &subject_repository.SubjectRepository{})
	subject_handler := subject_handler.NewSubjectHandler(subject_service, log)

	// init route
	router.Post("/subjects", subject_handler.CreateSubject)
	router.Put("/subjects", subject_handler.UpdateSubject)
	router.Delete("/subjects/{id}", subject_handler.DeleteSubject)

	// init server
	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}
	// start server
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
