package app

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/jure_assignments_handler"
	"main/internal/lib/liblogger"
	"main/internal/lib/supportRequest"
	"main/internal/middleware/midlogger"
	"main/internal/rabbitmq/consumer"
	"main/internal/repositories/jury_assignments_repository"
	"main/internal/services/jury_assignments_service"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	server *http.Server
	log    *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {

	// Init storage
	storage := postgresql.MustPosgreSQL(cfg.GetDataSourceName())
	log.Info("storage are enabled")

	// init chi router
	router := chi.NewRouter()

	app := &App{log: log}

	// init cors
	app.initCors(router, cfg.AdditionalAddressesConfig)
	// init middlewares
	router.Use(midlogger.New(log))
	router.Use(middleware.URLFormat)

	// init support for request
	supportReq := supportRequest.New(&cfg.AdditionalAddressesConfig)

	// init orm
	gormOrm := orm.NewGormORM(storage)
	// init repository
	repository := jury_assignments_repository.NewJuryAssignmentsRepository()
	// init sevice
	service := jury_assignments_service.NewJuryAssignmentsService(log, gormOrm, supportReq, repository)
	// init handler
	handler := jure_assignments_handler.NewJureAssignmentHandler(service, log)

	// init routes
	app.initRoutes(router, handler)

	// init rabbit
	rabbitConsumer := consumer.New(log, cfg.AddressRabbitPath, service)
	rabbitConsumer.Start(context.TODO(), cfg.QueueName)

	app.server = &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}

	return app
}

func (a *App) initCors(router *chi.Mux, cfg config.AdditionalAddressesConfig) {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{cfg.ReactVision}, // React URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // В секундах
	}
	router.Use(cors.Handler(corsOptions))
}

func (a *App) initRoutes(router *chi.Mux, handler *jure_assignments_handler.JuryAssignmentHandler) {
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Get("/api/jury-assignments", handler.GetAllJuryAssignments)
	router.Get("/api/jury-assignments/{id}", handler.GetJuryAssignmentsByID)
	router.Get("/api/jury-assignments/event/{id}", handler.GetAllByEventId)
	router.Get("/api/jury-assignments/jury/{id}", handler.GetAllByJuryId)

	router.Post("/api/jury-assignments", handler.CreateJuryAssignments)
	router.Post("/api/jury-assignments/many", handler.CreateManyJuryAssignments)

	router.Put("/api/jure-assignments/{id}", handler.UpdateJuryAssignments)

	router.Delete("/api/jure-assignments/{id}", handler.DeleteJuryAssignments)
	router.Post("/api/jury-assignments/delete/many", handler.DeleteManyAssigments)
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
		a.log.Error("failed to stop server", liblogger.Err(err))
		return
	}
	a.log.Info("server stopped")
}
