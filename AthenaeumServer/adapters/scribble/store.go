package scribble

import (
	"os"
	"strings"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/sdomino/scribble"
)

type ScribbleStore struct {
	*scribble.Driver
	logger logging.Logger
}

var bookCollection string = "books"

func initBookCollection(dbLocation string, logger logging.Logger) {
	collectionLocation := collectionLocation(dbLocation, bookCollection)
	logger.Debugf("Initalising collection %s at location %s", bookCollection, collectionLocation)
	//ignore error if cannot creat collection location
	_ = os.Mkdir(collectionLocation, os.ModePerm)
}

func NewScribbleStore(dbLocation string, logger logging.Logger) (*ScribbleStore, error) {
	logger.Debugf("Creating new ScribleStore")
	db, err := scribble.New(dbLocation, nil)
	logger.Debugf("Database client created with location %s", dbLocation)
	initBookCollection(dbLocation, logger)
	return &ScribbleStore{db, logger}, err
}

func collectionLocation(dbLocation string, collection string) string {
	if strings.HasSuffix(dbLocation, "/") {
		return dbLocation + collection
	} else {
		return dbLocation + "/" + collection
	}

}
