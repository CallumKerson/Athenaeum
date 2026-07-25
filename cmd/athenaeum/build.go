package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/CallumKerson/Athenaeum/internal/feed"
	"github.com/CallumKerson/Athenaeum/internal/scan"
	"github.com/CallumKerson/Athenaeum/internal/site"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// mediaPath is the URL prefix the media root is served under. It must keep
// matching the proxy's configuration: it forms part of every item's GUID.
const mediaPath = "/media"

var (
	errNoMediaRoot = errors.New("no media root configured: set Media.Root or pass --media-root")
	errNoSiteRoot  = errors.New("no output directory configured: set Site.Root or pass --out")
	errOvercast    = errors.New("overcast ping failed")
)

type buildFlags struct {
	configPath string
	mediaRoot  string
	siteRoot   string
	host       string
	cachePath  string
	noCache    bool
	noSweep    bool
	verbose    bool
}

func NewBuildCommand() *cobra.Command {
	var flags buildFlags

	cmd := &cobra.Command{
		Use:          "build",
		Short:        "Generate the podcast feed site from the audiobook library",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(cmd, &flags)
		},
	}

	cmd.Flags().StringVarP(&flags.configPath, "config", "c", "", "path to config file")
	cmd.Flags().StringVar(&flags.mediaRoot, "media-root", "", "directory holding the audiobook library")
	cmd.Flags().StringVarP(&flags.siteRoot, "out", "o", "", "directory to write the generated site to")
	cmd.Flags().StringVar(&flags.host, "host", "", "external base URL the site is served from")
	cmd.Flags().StringVar(&flags.cachePath, "cache", "", "path to the m4b duration cache")
	cmd.Flags().BoolVar(&flags.noCache, "no-cache", false, "re-read every m4b instead of using the cache")
	cmd.Flags().BoolVar(&flags.noSweep, "no-sweep", false, "keep files from previous builds that are now stale")
	cmd.Flags().BoolVarP(&flags.verbose, "verbose", "v", false, "log every file written")

	return cmd
}

func runBuild(cmd *cobra.Command, flags *buildFlags) error {
	logger := buildLogger(cmd.ErrOrStderr(), flags.verbose)

	cfg, err := resolveConfig(flags, cmd.ErrOrStderr())
	if err != nil {
		return err
	}
	excludedGenres, err := cfg.ExclusionsFromMainFeed.GetGenres()
	if err != nil {
		return err
	}

	start := time.Now()
	books, cache, err := scanLibrary(cfg, flags, logger)
	if err != nil {
		return err
	}
	logger.Info("scanned library",
		"books", len(books), "cached", cache.Hits, "parsed", cache.Misses, "took", time.Since(start).String())

	result, err := site.Build(cfg.Site.Root, site.Plan(books, excludedGenres), rendererFor(cfg), !flags.noSweep, logger)
	if err != nil {
		return err
	}
	logger.Info("built site",
		"root", cfg.Site.Root, "feeds", result.Feeds, "written", result.Written, "removed", result.Removed,
		"took", time.Since(start).String())

	if cfg.ThirdParty.NotifyOvercast && (result.Written > 0 || result.Removed > 0) {
		if err = notifyOvercast(cmd.Context(), cfg.Host); err != nil {
			logger.Warn("could not notify Overcast", "error", err)
		} else {
			logger.Info("notified Overcast", "urlprefix", cfg.Host)
		}
	}

	return nil
}

// resolveConfig layers the flags over the config file and checks that the two
// paths the build cannot invent for itself are present.
func resolveConfig(flags *buildFlags, out io.Writer) (*BuildConfig, error) {
	cfg, err := LoadBuildConfig(flags.configPath, out)
	if err != nil {
		return nil, err
	}
	if flags.mediaRoot != "" {
		cfg.Media.Root = flags.mediaRoot
	}
	if flags.siteRoot != "" {
		cfg.Site.Root = flags.siteRoot
	}
	if flags.host != "" {
		cfg.Host = flags.host
	}

	if cfg.Media.Root == "" {
		return nil, errNoMediaRoot
	}
	if cfg.Site.Root == "" {
		return nil, errNoSiteRoot
	}
	return cfg, nil
}

// scanLibrary walks the media root, returning the cache alongside the books so
// the caller can report how much of the scan the cache saved.
//
// Cache problems are never fatal: the worst case is re-reading every m4b.
func scanLibrary(
	cfg *BuildConfig,
	flags *buildFlags,
	logger *slog.Logger,
) ([]audiobooks.Audiobook, *scan.Cache, error) {
	cachePath := flags.cachePath
	if cachePath == "" {
		var err error
		if cachePath, err = DefaultCachePath(); err != nil {
			return nil, nil, err
		}
	}

	cache := scan.NewCache()
	if !flags.noCache {
		loaded, err := scan.LoadCache(cachePath)
		if err != nil {
			logger.Warn("could not read cache, re-reading every m4b", "path", cachePath, "error", err)
		}
		cache = loaded
	}

	logger.Info("scanning library", "root", cfg.Media.Root)
	books, err := scan.Library(cfg.Media.Root, cache, logger)
	if err != nil {
		return nil, nil, err
	}

	if !flags.noCache {
		if err = cache.Save(cachePath); err != nil {
			logger.Warn("could not save cache", "path", cachePath, "error", err)
		}
	}
	return books, cache, nil
}

func rendererFor(cfg *BuildConfig) *feed.Renderer {
	return &feed.Renderer{
		Host:               cfg.Host,
		MediaPath:          mediaPath,
		ImageLink:          fmt.Sprintf("%s/static/itunes_image.jpg", cfg.Host),
		Explicit:           cfg.Podcast.Explicit,
		Language:           cfg.Podcast.Language,
		Author:             cfg.Podcast.Author,
		Email:              cfg.Podcast.Email,
		HandlePreUnixEpoch: cfg.Podcast.PreUnixEpoch.Handle,
	}
}

func buildLogger(out io.Writer, verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: level}))
}

// notifyOvercast asks Overcast to re-fetch every feed under the host, so
// subscribers see new books without waiting for their next poll.
func notifyOvercast(ctx context.Context, host string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	endpoint := "https://overcast.fm/ping?urlprefix=" + url.QueryEscape(host)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w with status %s", errOvercast, response.Status)
	}
	return nil
}
