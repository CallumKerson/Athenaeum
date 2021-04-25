package main

import (
	"os"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/actions/books"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/actions/mp4"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/actions/podcasts"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/filehash"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/podcast"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/rest"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/scribble"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/tag"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/web"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/xid"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

func main() {
	cfg := loadConfig()
	logger := logging.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)
	idAdapter := xid.NewAdapter(logger)
	db, err := scribble.NewScribbleStore(cfg.DBLocation, logger)
	if err != nil {
		logger.Fatalf("Cannot create a db client: %s", err)
	}
	tag := tag.NewTagProvider(logger)
	mp4 := mp4.NewMP4DurationProvider(logger)
	fileHasher := filehash.NewMD5FileHasher()

	bookFactory := books.NewBookFactory(idAdapter, db, cfg.LibraryHome(), tag, mp4, fileHasher, logger)
	bookRetriever := books.NewRetriever(db, logger)
	bookScanner := books.NewScanner(logger, cfg.libraryHome, *bookFactory, *bookRetriever, fileHasher)
	podcastFactory := podcasts.NewPodcastFactory(cfg.LibraryHome(), logger)
	podcastGenerator := podcast.NewPodcastGenerator()
	podcastSerializer := podcasts.NewPodcastSerializer(podcastGenerator, logger)

	_, err = bookScanner.Scan()
	if err != nil {
		logger.Errorf("Cannot Scan for new files: %s", err)
	}

	webHandler, err := web.New(&cfg, bookRetriever, podcastFactory, podcastSerializer, logger)
	if err != nil {
		logger.Fatalf("failed to setup web handler: %v", err)
	}
	restHandler := rest.New(logger, bookRetriever, bookFactory, bookScanner)

	setupServer(cfg, webHandler, restHandler, logger)

}
