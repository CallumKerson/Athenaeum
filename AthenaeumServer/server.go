package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/middlewares"
	"github.com/gorilla/mux"
)

func setupServer(config config, web http.Handler, rest http.Handler, logger logging.Logger) {

	rest = middlewares.WithBasicAuth(logger, rest,
		middlewares.UserVerifierFunc(func(ctx context.Context, name, secret string) bool {
			return name == config.APIUser && secret == config.APIPassword
		}),
	)

	router := mux.NewRouter()
	logger.Debugf("Adding api handler")
	router.PathPrefix("/api").Handler(http.StripPrefix("/api", rest))
	logger.Debugf("Adding web handler")
	router.PathPrefix("/").Handler(web)

	handler := middlewares.WithRequestLogging(logger, router)

	srv := &http.Server{
		Handler:      handler,
		Addr:         config.Address,
		ReadTimeout:  config.GracefulTimeout,
		WriteTimeout: config.GracefulTimeout,
	}

	// Start Server
	go func() {
		logger.Infof("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("Cannot start server", err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv, logger)
}

func waitForShutdown(srv *http.Server, logger logging.Logger) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_ = srv.Shutdown(ctx)

	logger.Infof("Shutting down")
	os.Exit(0)
}
