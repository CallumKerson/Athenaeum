package main

import (
	"os"

	"github.com/CallumKerson/loggerrific"

	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
)

func Run(logger loggerrific.Logger) error {
	logger.Infoln("Setting up Server")

	return nil
}

func main() {
	logger := logrus.NewLogger()
	logger.SetLevelDebug()

	if err := Run(logger); err != nil {
		logger.WithError(err).Errorln("Error starting up server")
		os.Exit(1)
	}
}
