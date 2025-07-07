package main

import (
	"context"
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/event_handler"
	"main/internal/handlers/subject_handler"
	"main/internal/lib/liblogger"
	"main/internal/middleware/auth"
	"main/internal/middleware/midlogger"
	"main/internal/models/subject"
	"main/internal/repositories/event_repository"
	"main/internal/services/event_service"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"main/support/userrole"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const LocalFilePath = "config-yaml/dev.yaml"

var DockerFilePath string = os.Getenv("CONFIG_PATH")

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
	// init  middlewares cors
	corsOptions := cors.Options{
		AllowedOrigins: []string{cfg.ReactVision, cfg.JureAssignmentsService},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}
	router.Use(cors.Handler(corsOptions))
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)
	// add Authentication with JWT token
	router.Use(func(next http.Handler) http.Handler {
		return auth.AuthenticateMiddleware(next, cfg.Key)
	})

	// init subject service and handler
	subjectStorage := subject.NewSubjectsStorage()
	subjectHandler := subject_handler.NewSubjectHandler(subjectStorage, log)
	// init orm
	gormORM := orm.NewGormORM(storage)

	// init events service and handler
	eventService := event_service.NewEventService(gormORM, &event_repository.EventRepository{}, log)
	eventHandler := event_handler.NewEventHandler(eventService)

	// init subjects route
	router.Get("/events/subjects", subjectHandler.GetAllSubjects)

	// init events route
	router.With(auth.RoleBasedAccess(userrole.OrganizerRole)).Group(func(r chi.Router) {
		r.Post("/events", eventHandler.CreateEvent)
		r.Put("/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/events/{id}", eventHandler.DeleteEvent)
	})

	// router.Post("/events/byjson", eventHandler.CreateEventsByJSON)
	router.Post("/events/details/one", eventHandler.GetEventByFilterAndFields)
	router.Post("/events/details", eventHandler.GetEventsByFilterAndFields)
	router.Get("/events", eventHandler.GetAllEvents)
	router.Get("/events/{id}", eventHandler.GetEventByID)
	router.Get("/events/regional-stage", eventHandler.GetEventsTypeRegionalStage)
	router.Get("/events/stages/{id}", eventHandler.GetEventsTypeStageAndHisChilds)
	router.Get("/events/child/{id}", eventHandler.GetEventsByPreviousID)
	router.Post("/events/list", eventHandler.GetEventsByListID)

	// init server
	server := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}
	// init gorutine for start server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()
	log.Info("server started")
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", liblogger.Err(err))
		return
	}
	log.Info("server stopped")
}
