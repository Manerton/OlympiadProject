package app

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/event_handler"
	"main/internal/handlers/subject_handler"
	"main/internal/lib/liblogger"
	"main/internal/middleware/auth"
	"main/internal/middleware/midlogger"
	"main/internal/models/subject"
	"main/internal/repositories/event_repository"
	"main/internal/repositories/outbox_repository"
	"main/internal/services/event_service"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"main/rabbitmq"
	"main/rabbitmq/consumer"
	"main/rabbitmq/producer"
	"main/support/userrole"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	server *http.Server
	log    *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {

	// init storage
	storage := postgresql.MustPosgreSQL(cfg.GetDataSourceName())
	log.Info("storage are enabled")

	// init router
	router := chi.NewRouter()

	app := &App{log: log}
	// init cors
	app.initCors(router, cfg.AdditionalAddressesConfig)
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)
	// add Authentication with JWT token

	// router.Use(func(next http.Handler) http.Handler {
	// 	return auth.AuthenticateMiddleware(next, cfg.Key)
	// })

	// init subject service and handler
	subjectStorage := subject.NewSubjectsStorage()
	subjectHandler := subject_handler.NewSubjectHandler(subjectStorage, log)
	// init orm
	gormORM := orm.NewGormORM(storage)

	// init events service and handler
	eventService := event_service.NewEventService(gormORM, cfg.Services, &event_repository.EventRepository{}, &outbox_repository.OutboxRepository{}, log)
	eventHandler := event_handler.NewEventHandler(eventService)

	// init connection manager for rabbit
	connectionManager := rabbitmq.New(cfg.AddressRabbitPath, log)
	connectionManager.Start(context.TODO())

	// init rabbit consumer
	rabbitConsumer := consumer.New(log, connectionManager, eventService)
	rabbitConsumer.Start(context.TODO(), cfg.QueueName)

	// init rabbit producer
	rabbitProducer := producer.New(log, connectionManager, gormORM, &outbox_repository.OutboxRepository{})
	rabbitProducer.Start(context.TODO())

	// init routes
	app.initRoutes(router, eventHandler, subjectHandler)

	// init server
	app.server = &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}

	return app
}

func (a *App) initRoutes(router *chi.Mux,
	eventHandler *event_handler.EventHandler,
	subjectHandler *subject_handler.SubjectHandler,
) {

	// init subjects route
	router.Get("/api/events/subjects", subjectHandler.GetAllSubjects)

	// init events route
	router.With(auth.RoleBasedAccess(userrole.AdminRole)).Group(func(r chi.Router) {
		r.Post("/api/events", eventHandler.CreateEvent)
		r.Put("/api/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/api/events/{id}", eventHandler.DeleteEvent)
	})

	router.Post("/api/events/details/one", eventHandler.GetEventByFilterAndFields)
	router.Post("/api/events/details", eventHandler.GetEventsByFilterAndFields)
	router.Post("/api/events/list", eventHandler.GetEventsByListID)

	router.Get("/api/events", eventHandler.GetAllEvents)
	router.Get("/api/events/{id}", eventHandler.GetEventByID)
	router.Get("/api/events/regional-stage", eventHandler.GetEventsTypeRegionalStage)
	router.Get("/api/events/class", eventHandler.GetEventsByClassType)
	router.Get("/api/events/stages/{id}", eventHandler.GetEventsTypeStageAndHisChilds)
	router.Get("/api/events/child/{id}", eventHandler.GetEventsByPreviousID)
}

func (a *App) initCors(router *chi.Mux, cfg config.AdditionalAddressesConfig) {
	corsOptions := cors.Options{
		AllowedOrigins: []string{cfg.ReactVision, cfg.JureAssignmentsService, cfg.ApiGateway},
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
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	a.log.Info("server starting")
	if err := a.server.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("failee to stop sever", liblogger.Err(err))
		return
	}
	a.log.Info("server stopped")
}
