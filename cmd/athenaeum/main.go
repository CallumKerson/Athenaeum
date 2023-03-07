package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/podcasts"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	defaultPort = 8080
)

type DummyService struct {
}

func (s *DummyService) IsReady(ctx context.Context) (bool, error) {
	return true, nil
}

func Run(port int, logger loggerrific.Logger) error {
	opts := &podcasts.FeedOpts{
		Title:       "Audiobooks",
		Description: "Like movies for your mind!",
		Explicit:    true,
		Language:    "EN",
		Author:      "A Person",
		Email:       "person@domain.test",
		Copyright:   "None",
	}
	podcastService := podcasts.NewService(opts, logger)
	httpHandler := transportHttp.NewHandler(podcastService, logger)

	logger.Debugln("Setting up Server")
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           httpHandler,
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

func main() {
	logger := logrus.NewLogger()

	if err := Run(defaultPort, logger); err != nil {
		logger.WithError(err).Errorln("Error starting up server")
		os.Exit(1)
	}
}
