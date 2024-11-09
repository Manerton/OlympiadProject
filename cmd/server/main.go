package main

import (
	"context"
	"main/internal/config"
	jureAssignmentHandler "main/internal/handlers/jureAssignmentsHandler"
	"main/internal/lib/liblogger"
	"main/internal/lib/midlogger"
	"main/internal/repositories/juryAssignmentsRepository"
	"main/internal/services/juryAssignmentsService"
	"main/internal/storage/postgresql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const localConfigPath = "config-yaml/local.yaml"

func main() {
	cfg := config.GetConfig(localConfigPath)

	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start jure-assignments-service")
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

	// init repository
	repository := juryAssignmentsRepository.NewJuryAssignmentsRepository()
	// init sevice
	service := juryAssignmentsService.NewJuryAssignmentsService(storage, repository)
	// init handler
	handler := jureAssignmentHandler.NewJureAssignmentHandler(service, log)

	// init route
	router.Get("/jure-assignments", handler.GetAllJuryAssignments)
	router.Get("/jure-assignments/{id}", handler.GetJuryAssignmentsByID)
	router.Get("/jure-assignments/jury/{event_id}", handler.GetAllJuryIDByEventID)
	router.Post("/jure-assignments", handler.CreateJuryAssignments)

	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// gorutine
	go func() {
		// start server
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()
	log.Info("server started")
	<-done

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", liblogger.Err(err))
		return
	}

	log.Info("server stopped")
}
