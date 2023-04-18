package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/CallumKerson/Athenaeum/internal/adapters/alfgmp4"
	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	audiobooksService "github.com/CallumKerson/Athenaeum/internal/audiobooks/service"
	mediaService "github.com/CallumKerson/Athenaeum/internal/media/service"
	"github.com/CallumKerson/Athenaeum/internal/memcache"
	overcastNotifier "github.com/CallumKerson/Athenaeum/internal/overcast/notifier"
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
			athenaeumClient := client.New((&cfg).GetClientOpts()...)
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
	logger := cfg.GetLogger()
	m4bMetadataReader := alfgmp4.NewMetadataReader()
	mediaSvc := mediaService.New(m4bMetadataReader, logger, cfg.GetMediaServiceOpts()...)
	boltAudiobookStore, err := bolt.NewAudiobookStore(logger, true, cfg.GetBoltDBOps()...)
	if err != nil {
		return err
	}
	var updaters []audiobooksService.ThirdPartyNotifier
	if cfg.ThirdParty.UpdateOvercast || cfg.ThirdParty.NotifyOvercast {
		updaters = append(updaters, overcastNotifier.New(cfg.Host, logger))
	}
	audiobookSvc := audiobooksService.New(mediaSvc, boltAudiobookStore, logger, updaters...)
	if errScan := audiobookSvc.UpdateAudiobooks(context.Background()); errScan != nil {
		return errScan
	}
	podcastSvc := podcastService.New(audiobookSvc, logger, cfg.GetPodcastServiceOpts()...)

	var httpHandler *transportHttp.Handler
	if cfg.Cache.Enabled {
		httpHandler = transportHttp.NewHandler(
			podcastSvc,
			cfg.Media.Root,
			transportHttp.WithCacheStore(memcache.NewStore(cfg.GetMemcacheOpts()...)),
			transportHttp.WithLogger(logger),
			transportHttp.WithVersion(Version),
		)
	} else {
		httpHandler = transportHttp.NewHandler(
			podcastSvc,
			cfg.Media.Root,
			transportHttp.WithLogger(logger),
			transportHttp.WithVersion(Version),
		)
	}

	return transportHttp.Serve(httpHandler, cfg.Port, logger)
}
