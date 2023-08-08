package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/CallumKerson/Athenaeum/internal/adapters/alfgmp4"
	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	audiobooksService "github.com/CallumKerson/Athenaeum/internal/audiobooks/service"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	"github.com/CallumKerson/Athenaeum/internal/memcache"
	overcastNotifier "github.com/CallumKerson/Athenaeum/internal/notifiers/overcast"
	podcastService "github.com/CallumKerson/Athenaeum/internal/podcasts/service"
	transportHttp "github.com/CallumKerson/Athenaeum/internal/transport/http"
	"github.com/CallumKerson/Athenaeum/pkg/client"
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
	rootCmd.AddCommand(NewUpdateCommand())
	return rootCmd
}

func NewRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "run",
		Short:        "Run the server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(&cfg)
		},
	}
}

func NewUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "update",
		Short:        "Triggers an update on the running athenaeum instance.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			athenaeumClient := client.New((&cfg).Host, client.WithVersion(Version))
			err := athenaeumClient.Update(context.Background())
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Updated Athenaeum at %s\n", cfg.Host)
			return err
		},
	}
}

func runServer(cfg *Config) error {
	logger := getLogger(cfg)

	m4bMetadataReader := alfgmp4.NewMetadataReader()
	mediaSvc := mediaService.New(m4bMetadataReader, cfg.Media.Root, mediaService.WithLogger(logger))
	boltAudiobookStore, err := bolt.NewAudiobookStore(cfg.DB.Root, bolt.WithLogger(logger))
	if err != nil {
		return err
	}
	audiobookOpts := []audiobooksService.Option{audiobooksService.WithLogger(logger)}
	if cfg.ThirdParty.UpdateOvercast || cfg.ThirdParty.NotifyOvercast {
		audiobookOpts = append(audiobookOpts,
			audiobooksService.WithThirdPartyNotifier(overcastNotifier.New(cfg.Host, overcastNotifier.WithLogger(logger))))
	}
	if len(cfg.ExcludsionsFromMainFeed.Genres) > 0 {
		genres, err := cfg.ExcludsionsFromMainFeed.GetGenres()
		if err != nil {
			return err
		}
		audiobookOpts = append(audiobookOpts, audiobooksService.WithGenresToExludeFromAllAudiobooks(genres...))
	}
	audiobookSvc := audiobooksService.New(mediaSvc, boltAudiobookStore, audiobookOpts...)
	if errScan := audiobookSvc.UpdateAudiobooks(context.Background()); errScan != nil {
		return errScan
	}
	podcastSvc := podcastService.New(audiobookSvc, cfg.Host, transportHttp.MediaPath, podcastService.WithLogger(logger),
		podcastService.WithPodcastFeedInfo(
			cfg.Podcast.Explicit,
			cfg.Podcast.Language,
			cfg.Podcast.Author,
			cfg.Podcast.Email,
			cfg.Podcast.Copyright,
			fmt.Sprintf("%s%s/itunes_image.jpg", cfg.Host, transportHttp.StaticPath),
		),
		podcastService.WithHandlePreUnixEpoch(cfg.Podcast.PreUnixEpoch.Handle))

	httpHandlerOpts := []transportHttp.HandlerOption{transportHttp.WithLogger(logger), transportHttp.WithVersion(Version)}
	if cfg.Cache.Enabled {
		httpHandlerOpts = append(httpHandlerOpts, transportHttp.WithCacheStore(
			memcache.NewStore(
				memcache.WithTTL(cfg.Cache.GetTTL()),
				memcache.WithCapacity(cfg.Cache.Length),
			),
		))
	}

	httpHandler := transportHttp.NewHandler(
		podcastSvc,
		cfg.Media.Root,
		httpHandlerOpts...,
	)
	return transportHttp.Serve(httpHandler, cfg.Port, logger)
}

func getLogger(cfg *Config) *logrus.Logger {
	log := logrus.NewLogger()
	level := strings.ToLower(cfg.Log.Level)
	switch level {
	case "debug":
		log.SetLevelDebug()
	case "warn":
		log.SetLevelWarn()
	case "error":
		log.SetLevelError()
	default:
		log.SetLevelInfo()
	}
	return log
}
