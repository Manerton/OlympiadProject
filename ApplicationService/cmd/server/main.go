package main

import (
	"context"
	"log/slog"
	"main/internal/config"
	ApplicationHandler "main/internal/handlers"
	"main/internal/lib/liblogger"
	"main/internal/middleware/auth"
	"main/internal/middleware/midlogger"
	rabbit "main/internal/rabbitmq"
	"main/internal/rabbitmq/consumer"
	ApplicationRepository "main/internal/repositories"
	ApplicationService "main/internal/services"
	"main/internal/storage/orm"
	"main/internal/storage/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const LocalFilePath = "../../config-yaml/docker.yaml"

func main() {

	// Init config
	cfg := config.GetConfig(LocalFilePath)

	// Init logger
	log := liblogger.SetupLogger(cfg.Env)
	log.Info("start application-service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	// storage, err := postgresql.NewPostgreSQL(cfg.GetDataSourceName())
	// if err != nil {
	// 	log.Error("failed to init storage", liblogger.Err(err))
	// } else {
	// 	log.Info("storage is enabled")
	// }

	storage := postgresql.MustPosgreSQL(cfg.GetDataSourceName())

	// init orm
	gormORM := orm.NewGormORM(storage)

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
	// router.Use(func(next http.Handler) http.Handler {
	// 	return auth.AuthenticateMiddleware(next, cfg.Key)
	// })

	applicationsRepo := &ApplicationRepository.ApplicationRepository{}

	// init application service and handler
	applicationService := ApplicationService.NewApplicationService(gormORM, applicationsRepo)
	applicationHandler := ApplicationHandler.NewApplicationHandler(applicationService, log)

	router.Get("/applications/{id}", applicationHandler.GetApplicationByID) //wtf

	router.Get("/applications", applicationHandler.GetAllApplications)

	router.Get("/applicationsCount", applicationHandler.GetCountApplications)

	router.Get("/applicationsByFilter", applicationHandler.GetByFilter)

	router.Get("/applications/user/{userID}", applicationHandler.GetApplicationsByUserID)

	router.Get("/applications/event/{eventID}", applicationHandler.GetApplicationsByEventID)

	router.Get("/applications/school/{schoolID}", applicationHandler.GetApplicationsBySchoolID)

	router.Get("/applications/school-list/", applicationHandler.GetApplicationsBySchoolListID)

	router.Post("/applications/create", applicationHandler.CreateApplication)

	router.Put("/applications/{id}", applicationHandler.UpdateApplicationStatus)
	router.Delete("/applications/{id}", applicationHandler.DeleteApplication)

	// init applications route
	router.With(auth.RoleBasedAccess("4")).Group(func(r chi.Router) {

	})

	router.With(auth.RoleBasedAccess("2", "4")).Group(func(r chi.Router) {

	})

	router.With(auth.RoleBasedAccess("3")).Group(func(r chi.Router) {

	})

	router.With(auth.RoleBasedAccess("2", "4")).Group(func(r chi.Router) {

	})

	router.With(auth.RoleBasedAccess("3", "4")).Group(func(r chi.Router) {

	})

	//init rabbitMQ

	// init connection manager for rabbit
	connectionManager := rabbit.New(cfg.AddressRabbitPath, log)
	connectionManager.Start(context.TODO())

	// init rabbit consumer
	rabbitConsumer := consumer.New(log, connectionManager, applicationService)
	rabbitConsumer.Start(context.TODO(), cfg.QueueName)
	log.Info("rabbit started")

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
