package main

import (
	"os"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/config"
	"github.com/CallumKerson/Athenaeum/internal/podcasts"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	defaultPort = 8080
)

func Run(port int, logger loggerrific.Logger) error {
	cfg, err := config.New(logger)
	if err != nil {
		return err
	}
	podcastService := podcasts.NewService(&cfg.Podcast, logger)
	httpHandler := transportHttp.NewHandler(podcastService, logger)

	return Serve(httpHandler, port, logger)
}

func main() {
	logger := logrus.NewLogger()

	if err := Run(defaultPort, logger); err != nil {
		logger.WithError(err).Errorln("Error starting up server")
		os.Exit(1)
	}
}
