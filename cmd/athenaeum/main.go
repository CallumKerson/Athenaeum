package main

import (
	"os"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	"github.com/CallumKerson/Athenaeum/internal/podcasts"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	defaultPort = 8080
)

func Run(port int, logger loggerrific.Logger) error {
	cfg, err := NewConfig(port, logger)
	if err != nil {
		return err
	}
	mediaSvc := mediaService.New(cfg.Media.Root, logger)
	podcastService := podcasts.NewService(cfg.GetMediaHost(), &cfg.Podcast.Opts, mediaSvc, logger)
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
