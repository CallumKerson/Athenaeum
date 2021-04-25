package xid

import (
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
	"github.com/rs/xid"
)

func NewAdapter(logger logging.Logger) *Adapter {
	logger.Debugf("Creating new XidAdapter")
	return &Adapter{logger: logger}
}

type Adapter struct {
	logger logging.Logger
}

func (xidAdapter *Adapter) Get() string {
	return xid.New().String()
}
