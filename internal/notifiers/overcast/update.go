package overcast

import (
	"context"
	"fmt"
	"time"

	"github.com/CallumKerson/loggerrific"
	noOpLogger "github.com/CallumKerson/loggerrific/noop"
	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"
)

const (
	OvercastHost = "https://overcast.fm"
)

type Notifier struct {
	host      string
	urlPrefix string
	logger    loggerrific.Logger
}

func New(urlPrefix string, opts ...Option) *Notifier {
	notifier := &Notifier{
		host:      OvercastHost,
		urlPrefix: urlPrefix,
		logger:    noOpLogger.New(),
	}
	for _, opt := range opts {
		opt(notifier)
	}
	notifier.logger.Infoln("Will update", notifier, "for URL prefix", urlPrefix)
	return notifier
}

func (n *Notifier) Notify(ctx context.Context) error {
	err := requests.
		URL(n.host).
		Client(httpretry.NewDefaultClient(
			httpretry.WithMaxRetryCount(8),
			httpretry.WithBackoffPolicy(
				httpretry.ExponentialBackoff(3*time.Minute, 120*time.Minute, 0),
			),
		)).
		Param("urlprefix", n.urlPrefix).
		Path("ping").
		Fetch(ctx)
	if err != nil {
		return err
	}
	n.logger.Infoln("Notified", n, "with urlprefix", n.urlPrefix)
	return nil
}

func (n *Notifier) String() string {
	return fmt.Sprintf("Overcast (%s)", n.host)
}
