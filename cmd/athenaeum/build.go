package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/CallumKerson/Athenaeum/internal/cover"
	"github.com/CallumKerson/Athenaeum/internal/feed"
	"github.com/CallumKerson/Athenaeum/internal/fsutil"
	"github.com/CallumKerson/Athenaeum/internal/scan"
	"github.com/CallumKerson/Athenaeum/internal/site"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// overcastPendingName sits beside the duration cache while an Overcast ping is
// owed, because the build after a failed ping may have nothing new to write.
const overcastPendingName = "overcast-ping-pending"

// coverCacheName is the directory beside the duration cache that holds cover
// thumbnails, so a rebuild only decodes covers that are new or changed.
const coverCacheName = "covers"

// overcastPingURL is a variable only so that tests can point it at a local server.
var overcastPingURL = "https://overcast.fm/ping"

var (
	errNoMediaRoot = errors.New("no media root configured: set Media.Root or pass --media-root")
	errNoSiteRoot  = errors.New("no output directory configured: set Site.Root or pass --out")
	errNoHost      = errors.New("no host configured: set Host or pass --host")
	errOvercast    = errors.New("overcast ping failed")
	errNoBooks     = errors.New("no audiobooks found in the media root: refusing to empty every feed " +
		"(pass --allow-empty if that is intended)")
)

type buildFlags struct {
	configPath string
	mediaRoot  string
	siteRoot   string
	host       string
	cachePath  string
	noCache    bool
	noSweep    bool
	allowEmpty bool
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
	cmd.Flags().BoolVar(&flags.allowEmpty, "allow-empty", false, "build even when the media root holds no audiobooks")
	cmd.Flags().BoolVarP(&flags.verbose, "verbose", "v", false, "log every file written")

	return cmd
}

func runBuild(cmd *cobra.Command, flags *buildFlags) error {
	logger := buildLogger(cmd.ErrOrStderr(), flags.verbose)

	cfg, err := resolveConfig(flags, cmd.ErrOrStderr())
	if err != nil {
		return err
	}
	// Built before the scan so that a bad host fails the build straight away.
	renderer, err := feed.New(cfg.Host, feed.Channel{
		Explicit:           cfg.Podcast.Explicit,
		Language:           cfg.Podcast.Language,
		Author:             cfg.Podcast.Author,
		Email:              cfg.Podcast.Email,
		HandlePreUnixEpoch: cfg.Podcast.PreUnixEpoch.Handle,
		Artwork:            true,
	})
	if err != nil {
		return err
	}
	excludedFromMainFeed, err := cfg.ExclusionsFromMainFeed.GetGenres()
	if err != nil {
		return err
	}
	excludedFromSite, err := cfg.ExclusionsFromSite.GetGenres()
	if err != nil {
		return err
	}
	exclusions := site.Exclusions{MainFeed: excludedFromMainFeed, Site: excludedFromSite}

	start := time.Now()
	books, cache, err := scanLibrary(cfg, flags, logger)
	if err != nil {
		return err
	}
	logger.Info("scanned library",
		"books", len(books), "cached", cache.Hits, "parsed", cache.Misses, "took", time.Since(start).String())

	covers, err := coverSource(cfg, flags)
	if err != nil {
		return err
	}
	result, err := site.Build(cfg.Site.Root, site.Plan(books, exclusions), renderer,
		site.Options{Covers: covers, Sweep: !flags.noSweep}, logger)
	if err != nil {
		return err
	}
	if err := covers.Prune(); err != nil {
		logger.Warn("could not prune the cover cache", "error", err)
	}
	logger.Info("built site",
		"root", cfg.Site.Root, "feeds", result.Feeds, "written", result.Written, "removed", result.Removed,
		"took", time.Since(start).String())

	if cfg.ThirdParty.NotifyOvercast {
		cachePath, err := resolveCachePath(flags)
		if err != nil {
			return err
		}
		pendingPath := filepath.Join(filepath.Dir(cachePath), overcastPendingName)
		pingOvercast(cmd.Context(), renderer.Host(), result.Written > 0 || result.Removed > 0, pendingPath, logger)
	}

	return nil
}

// coverSource reads covers from the media library, caching their thumbnails
// beside the duration cache unless --no-cache asks for everything afresh.
func coverSource(cfg *BuildConfig, flags *buildFlags) (*cover.Source, error) {
	if flags.noCache {
		return cover.NewSource(cfg.Media.Root, ""), nil
	}
	cachePath, err := resolveCachePath(flags)
	if err != nil {
		return nil, err
	}
	return cover.NewSource(cfg.Media.Root, filepath.Join(filepath.Dir(cachePath), coverCacheName)), nil
}

// resolveConfig layers the flags over the config file and checks that the
// paths and host the build cannot invent for itself are present.
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
	if cfg.Media.Root, err = expandHome(cfg.Media.Root); err != nil {
		return nil, err
	}
	if cfg.Site.Root, err = expandHome(cfg.Site.Root); err != nil {
		return nil, err
	}
	if cfg.Host == "" {
		return nil, errNoHost
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
	cachePath, err := resolveCachePath(flags)
	if err != nil {
		return nil, nil, err
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
	// An unmounted or emptied media root would otherwise publish empty feeds,
	// sweep every other feed away and prune the whole cache. Checked before the
	// cache is saved so that the cache survives too.
	if len(books) == 0 && !flags.allowEmpty {
		return nil, nil, fmt.Errorf("%w: %s", errNoBooks, cfg.Media.Root)
	}

	if !flags.noCache {
		if err = cache.Save(cachePath); err != nil {
			logger.Warn("could not save cache", "path", cachePath, "error", err)
		}
	}
	return books, cache, nil
}

func resolveCachePath(flags *buildFlags) (string, error) {
	if flags.cachePath == "" {
		return DefaultCachePath()
	}
	return expandHome(flags.cachePath)
}

func buildLogger(out io.Writer, verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: level}))
}

// pingOvercast notifies Overcast when the build changed something, or when an
// earlier build's ping failed. A failure is never fatal: the feeds are already
// written, and subscribers will see them at their next poll regardless.
func pingOvercast(ctx context.Context, host string, changed bool, pendingPath string, logger *slog.Logger) {
	_, statErr := os.Stat(pendingPath)
	pending := statErr == nil
	if !changed && !pending {
		return
	}

	if err := notifyOvercast(ctx, host); err != nil {
		logger.Warn("could not notify Overcast, will retry on the next build", "error", err)
		if writeErr := fsutil.WriteAtomic(pendingPath, nil); writeErr != nil {
			logger.Warn("could not record the pending Overcast ping", "path", pendingPath, "error", writeErr)
		}
		return
	}
	logger.Info("notified Overcast", "urlprefix", host)
	if pending {
		if err := os.Remove(pendingPath); err != nil {
			logger.Warn("could not clear the pending Overcast ping", "path", pendingPath, "error", err)
		}
	}
}

// notifyOvercast asks Overcast to re-fetch every feed under the host, so
// subscribers see new books without waiting for their next poll.
func notifyOvercast(ctx context.Context, host string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	endpoint := overcastPingURL + "?urlprefix=" + url.QueryEscape(host)
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
