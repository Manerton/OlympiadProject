package app

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/handlers/auth_handler"
	"main/internal/handlers/user_handler"
	"main/internal/lib/jwttoken"
	"main/internal/lib/liblogger"
	"main/internal/repositories/participant_repository"
	"main/internal/repositories/user_repository"
	"main/internal/services/auth_service"
	"main/internal/services/user_service"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type App struct {
	storage interface{}
	server  *http.Server
	log     *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {

	// init storage
	storage := postgresql.MustPosgreSQL(cfg.GetDataSourceName())
	log.Info("storage are enabled")
	// init orm
	gormORM := orm.NewGormORM(storage)
	// init jwtManager
	jwtManager := jwttoken.NewJWTManager(cfg.Key, time.Duration(cfg.Duration))
	// init repositories
	userRepository := &user_repository.UserRepository{}
	// init services
	authService := auth_service.New(log, gormORM, jwtManager, userRepository, &participant_repository.ParticipantRepository{})
	userService := user_service.New(log, gormORM, userRepository)
	// init handlers
	userHandler := user_handler.New(userService)
	authHandler := auth_handler.New(authService)

	router := chi.NewRouter()

	app := &App{log: log}
	// init cors
	app.initCors(router, cfg.AdditionalAddressesConfig)

	// TODO! init routes
	router.Post("/users/login", authHandler.Login)
	router.Post("/users/register", authHandler.Register)
	router.Get("/users/{id}", userHandler.GetUserById)
	router.Get("/users/list", userHandler.GetUsersByListId)

	// init server
	app.server = &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: router,
	}

	return app
}

func (a *App) initCors(router *chi.Mux, cfg config.AdditionalAddressesConfig) {
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
