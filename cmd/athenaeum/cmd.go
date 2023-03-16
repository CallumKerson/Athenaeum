package main

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/CallumKerson/Athenaeum/internal/adapters/alfgmp4"
	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	audiobooksService "github.com/CallumKerson/Athenaeum/internal/audiobooks/service"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
)

const (
	shortHelp = "an audiobook server that provides a podcast feed"
)

var (
	pathToConfig = ""
	cfg          Config
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Define our command
	rootCmd := &cobra.Command{
		Use:          "athenaeum",
		Short:        shortHelp,
		SilenceUsage: true,
		Version:      Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return InitConfig(&cfg, pathToConfig, cmd.OutOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&cfg)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.PersistentFlags().StringVarP(&pathToConfig, "config", "c", pathToConfig, "path to config file")

	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewRunCommand())
	return rootCmd
}

func NewRunCommand() *cobra.Command {
	pathToConfig := ""
	var cfg Config
	return &cobra.Command{
		Use:          "run",
		Short:        "Run the server",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return InitConfig(&cfg, pathToConfig, cmd.OutOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&cfg)
		},
	}
}

func runServer(cfg *Config) error {
	logger := logrus.NewLogger()
	setLogLevel(logger, cfg.GetLogLevel())
	m4bMetadataReader := alfgmp4.NewMetadataReader()
	mediaSvc := mediaService.New(m4bMetadataReader, logger, cfg.GetMediaServiceOpts()...)
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

	return transportHttp.Serve(httpHandler, cfg.Port, logger)
}

func setLogLevel(logger *logrus.Logger, level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		logger.SetLevelDebug()
	case "warn":
		logger.SetLevelWarn()
	case "error":
		logger.SetLevelError()
	default:
		logger.SetLevelInfo()
	}
}
