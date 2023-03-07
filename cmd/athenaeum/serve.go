package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CallumKerson/loggerrific"
)

func Serve(handler http.Handler, port int, logger loggerrific.Logger) error {
	logger.Debugln("Setting up Server")
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		ReadHeaderTimeout: 15 * time.Second,
	}

	logger.Infoln("Starting server on port", port)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.WithError(err).Errorln("Server Error")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logger.WithError(err).Errorln("Problem shutting down server")
		return err
	}

	logger.Infoln("Shut down server gracefully")
	return nil
}
