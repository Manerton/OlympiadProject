package main

import (
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/event_handler"
	"main/internal/handlers/subject_handler"
	"main/internal/lib/liblogger"
	"main/internal/middleware/midlogger"
	"main/internal/models/subject"
	"main/internal/repositories/event_repository"
	"main/internal/services/event_service"
	"main/internal/storage/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const LocalFilePath = "config-yaml/local.yaml"
const DebugFilePath = "C:/Users/assba/source/repos/Event-service/config-yaml/local.yaml"

func main() {
	// Init config
	cfg := config.GetConfig(DebugFilePath)

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
	// subjectService := subject_service.NewSubjectService(storage, &subject_repository.SubjectRepository{})
	subjectStorage := subject.NewSubjectsStorage()
	subjectHandler := subject_handler.NewSubjectHandler(subjectStorage, log)
	// init events service and handler
	eventService := event_service.NewEventService(storage, &event_repository.EventRepository{})
	eventHandler := event_handler.NewEventHandler(eventService, log)

	// init subjects route
	router.Get("/subjects", subjectHandler.GetAllSubjects)
	// router.Get("/subjects/{id}", subjectHandler.GetSubjectByID)
	// router.Post("/subjects", subjectHandler.CreateSubject)
	// router.Put("/subjects", subjectHandler.UpdateSubject)
	// router.Delete("/subjects/{id}", subjectHandler.DeleteSubject)

	// init events route
	router.Get("/events", eventHandler.GetAllEvents)
	router.Get("/events/{id}", eventHandler.GetEventByID)
	router.Get("/events/regional-stage", eventHandler.GetAllEventsTypeRegionalStage)
	router.Get("/events/child/{id}", eventHandler.GetEventsByPreviousID)
	router.Post("/events", eventHandler.CreateEvent)
	router.Put("/events/{id}", eventHandler.UpdateEvent)
	router.Delete("/events/{id}", eventHandler.DeleteEvent)

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
