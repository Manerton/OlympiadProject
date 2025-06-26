package main

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/handlers/jure_assignments_handler"
	"main/internal/lib/liblogger"
	"main/internal/middleware/auth"
	"main/internal/middleware/midlogger"
	"main/internal/repositories/jury_assignments_repository"
	"main/internal/services/jury_assignments_service"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const localConfigPath = "config-yaml/local.yaml"
const DebugLocalFilePath = "C:/code_folder/OlympiadProject/Jure-assignments-service/config-yaml/local.yaml"

func main() {
	cfg := config.GetConfig(DebugLocalFilePath)

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
	// init  middlewares cors
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

	// init repository
	repository := jury_assignments_repository.NewJuryAssignmentsRepository()
	// init orm
	gormOrm := orm.NewGormORM(storage)
	// init sevice
	service := jury_assignments_service.NewJuryAssignmentsService(log, gormOrm, repository)
	// init handler
	handler := jure_assignments_handler.NewJureAssignmentHandler(service, log)

	// init route
	router.Get("/jury-assignments", handler.GetAllJuryAssignments)
	router.Get("/jury-assignments/{id}", handler.GetJuryAssignmentsByID)
	router.Get("/jury-assignments/jury/{event_id}", handler.GetAllJuryIDByEventID)
	router.Post("/jury-assignments", handler.CreateJuryAssignments)
	router.Post("/jury-assignments/many", handler.CreateManyAssignmentsByOneJury)
	router.Put("/jure-assignments/{id}", handler.UpdateJuryAssignments)
	router.Delete("/jure-assignments/{id}", handler.UpdateJuryAssignments)

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
	fmt.Println(<-done)

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", liblogger.Err(err))
		return
	}

	log.Info("server stopped")
}
