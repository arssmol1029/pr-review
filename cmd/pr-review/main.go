package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"pr-review/internal/config"
	"pr-review/internal/database/sqlite"
	"pr-review/internal/server/handlers"
	"pr-review/internal/service"
	"syscall"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.GetSlogLevel()}),
	)
	log.Info("Starting application",
		"env", cfg.Env,
	)

	ctx := context.Background()

	repository, err := setupDatabase(ctx, log, cfg.Database)
	if err != nil {
		log.Error("Failed to setup database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := repository.Close(); err != nil {
			log.Error("Failed to close database", "error", err)
			os.Exit(1)
		}
	}()

	userService := service.NewUserService(log, repository)
	teamService := service.NewTeamService(log, repository)
	prService := service.NewPRService(log, repository)
	statsService := service.NewStatsService(log, repository)

	router := SetupRouter(log, teamService, userService, prService, statsService)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := repository.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Info("Starting server...")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err.Error())
		}
	}()

	log.Info("Server started")

	<-done
	log.Info("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to stop server", slog.String("error", err.Error()))

		return
	}

	log.Info("Server stopped")
}

func SetupRouter(
	logger *slog.Logger,
	teamService handlers.TeamService,
	userService handlers.UserService,
	prService handlers.PRService,
	statsService handlers.StatsService,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			next.ServeHTTP(w, r)
		})
	})

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	teamHandler := handlers.NewTeamHandler(logger, teamService)
	userHandler := handlers.NewUserHandler(logger, userService)
	prHandler := handlers.NewPRHandler(logger, prService)
	statsHandler := handlers.NewStatsHandler(logger, statsService)

	router.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReview)
	})
	router.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.Add)
		r.Get("/get", teamHandler.Get)
	})
	router.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", prHandler.Create)
		r.Post("/merge", prHandler.Merge)
		r.Post("/reassign", prHandler.Reassign)
	})
	router.Route("/stats", func(r chi.Router) {
		r.Get("/user", statsHandler.User)
		r.Get("/team", statsHandler.Team)
		r.Get("/total", statsHandler.Total)
	})

	return router
}

func setupDatabase(ctx context.Context, log *slog.Logger, dbCfg config.DatabaseConfig) (*sqlite.SQLiteRepository, error) {
	log.Info("Initializing database",
		"path", dbCfg.Path,
		"init_timeout", dbCfg.InitTimeout,
		"ping_timeout", dbCfg.PingTimeout,
	)

	initCtx, cancel := context.WithTimeout(ctx, dbCfg.InitTimeout)
	defer cancel()

	repo, err := sqlite.New(initCtx, dbCfg.Path)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, dbCfg.PingTimeout)
	defer cancel()

	if err := repo.Ping(pingCtx); err != nil {
		return nil, err
	}

	log.Info("Database initialized successfully")
	return repo, nil
}
