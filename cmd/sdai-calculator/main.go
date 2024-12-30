package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sdai-calculator/internal/config"
	"sdai-calculator/internal/generated/server"
	httphandlers "sdai-calculator/internal/http-server/handlers"
	mwlogger "sdai-calculator/internal/http-server/middleware/logger"
	"sdai-calculator/internal/lib/logger/sl"
	"sdai-calculator/internal/storage/postgresql"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	openapimw "github.com/go-openapi/runtime/middleware"
	nethttpmw "github.com/oapi-codegen/nethttp-middleware"
)

const (
	envLocal      string = "local"
	envProduction string = "production"
)

type ValidationError struct {
	Message string `json:"message"`
}

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting sdai-calculator service")

	storage, err := postgresql.New(cfg.StoragePath)

	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	handlers := httphandlers.NewHandlers(logger, storage)

	swagger, err := server.GetSwagger()
	if err != nil {
		logger.Error("failed to get swagger", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Get("/swagger/json", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(swagger)
		if err != nil {
			logger.Error("failed to serve swagger.json", sl.Err(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
	router.Handle("/swagger", openapimw.SwaggerUI(openapimw.SwaggerUIOpts{
		Path:    "/swagger",
		SpecURL: "/swagger/json",
	}, nil))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		err := json.NewEncoder(w).Encode(ValidationError{
			Message: "Not found",
		})

		if err != nil {
			return
		}
	})

	validator := nethttpmw.OapiRequestValidatorWithOptions(
		swagger,
		&nethttpmw.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
					return nil
				},
				IncludeResponseStatus: true,
			},
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				err := json.NewEncoder(w).Encode(ValidationError{
					Message: message,
				})

				if err != nil {
					return
				}
			},
		},
	)

	handler := server.HandlerWithOptions(
		server.NewStrictHandler(handlers, nil),
		server.ChiServerOptions{
			BaseURL:    "",
			BaseRouter: router,
			Middlewares: []server.MiddlewareFunc{
				mwlogger.New(logger),
				middleware.Recoverer,
				middleware.RequestID,
				middleware.URLFormat,
				validator,
			},
		},
	)

	logger.Info("starting server", slog.String("addr", cfg.Address))

	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("failed to start server", sl.Err(err))
	}

	// fixme: graceful shutdown doesn't work

	logger.Info("server started", slog.String("addr", cfg.Address))

	gracefulShutdown(httpServer, logger)

	logger.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	logger = logger.With(slog.String("env", env))

	return logger
}

// fixme: create setupRouter func

func gracefulShutdown(server *http.Server, logger *slog.Logger) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logger.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", sl.Err(err))
	} else {
		logger.Info("server exited gracefully")
	}
}
