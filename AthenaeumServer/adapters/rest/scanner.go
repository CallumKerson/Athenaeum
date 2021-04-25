package rest

import (
	"net/http"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/gorilla/mux"
)

type scanner interface {
	Scan() (int, error)
	LastScanned() time.Time
	BooksFoundLastScan() int
}

type scannerController struct {
	logger  logging.Logger
	scanner scanner
}

func addScannerAPI(router *mux.Router, scanner scanner, logger logging.Logger) {
	sController := &scannerController{
		logger:  logger,
		scanner: scanner,
	}

	logger.Debugf("Adding handlers for scanning API")
	router.HandleFunc("/v1/scanner", sController.get).Methods(http.MethodGet)
	router.HandleFunc("/v1/scanner", sController.scan).Methods(http.MethodPut)
}

type scanDataHolder struct {
	BooksFoundLastScan int       `json:"booksFoundLastScan"`
	LastScanned        time.Time `json:"timeLastScanned"`
	Scanning           bool      `json:"scanning"`
}

func (s *scannerController) get(wr http.ResponseWriter, req *http.Request) {
	data := scanDataHolder{
		s.scanner.BooksFoundLastScan(),
		s.scanner.LastScanned(),
		false,
	}
	respond(wr, http.StatusOK, data)
}

func (s *scannerController) scan(wr http.ResponseWriter, req *http.Request) {
	parsedData := scanDataHolder{}
	if err := readRequest(req, &parsedData); err != nil {
		s.logger.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, err)
		return
	}
	if parsedData.Scanning {
		s.scanner.Scan()
	}

	data := scanDataHolder{
		s.scanner.BooksFoundLastScan(),
		s.scanner.LastScanned(),
		false,
	}
	respond(wr, http.StatusOK, data)
}
