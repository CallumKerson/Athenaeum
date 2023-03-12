package main

import (
	"context"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	audiobooksService "github.com/CallumKerson/Athenaeum/internal/audiobooks/service"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	defaultPort = 8080
)

var (
	Version = "development"
	Commit  = "development"
	Date    = "development"
)

func Run(port int, logger loggerrific.Logger) error {
	cfg, err := NewConfig(port, logger)
	if err != nil {
		return err
	}
	mediaSvc := mediaService.New(logger, cfg.GetMediaServiceOpts()...)
	boltAudiobookStore, err := bolt.NewAudiobookStore(logger, true, cfg.GetBoltDBOps()...)
	if err != nil {
		return err
	}
	audiobookSvc := audiobooksService.New(mediaSvc, boltAudiobookStore, logger)
	if errScan := audiobookSvc.UpdateAudiobooks(context.Background()); errScan != nil {
		return errScan
	}
	podcastSvc := podcastService.New(audiobookSvc, logger, cfg.GetPodcastServiceOpts()...)
	httpHandler := transportHttp.NewHandler(podcastSvc, audiobookSvc, logger, cfg.GetHTTPHandlerOpts()...)

	return Serve(httpHandler, port, logger)
}

func main() {
	showVersion := flag.BoolP("version", "v", false, "prints version and exits")
	flag.Parse()
	if *showVersion {
		fmt.Println("version: ", Version)
		fmt.Println("commit:  ", Commit)
		fmt.Println("built at:", Date)
		return
	}

	logger := logrus.NewLogger()

	if err := Run(defaultPort, logger); err != nil {
		logger.WithError(err).Errorln("Error starting up server")
		os.Exit(1)
	}
}
