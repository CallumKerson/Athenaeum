package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"

	"github.com/CallumKerson/loggerrific"
)

const (
	OvercastHost = "https://overcast.fm"
)

type Notifier struct {
	host      string
	urlPrefix string
	logger    loggerrific.Logger
}

func New(urlPrefix string, logger loggerrific.Logger) *Notifier {
	upd := &Notifier{host: OvercastHost, urlPrefix: urlPrefix, logger: logger}
	logger.Infoln("Will update", upd, "for URL prefix", urlPrefix)
	return upd
}

func (u *Notifier) Notify(ctx context.Context) error {
	err := requests.
		URL(u.host).
		Client(httpretry.NewDefaultClient(
			httpretry.WithMaxRetryCount(6),
			httpretry.WithBackoffPolicy(
				httpretry.ExponentialBackoff(1*time.Minute, 15*time.Minute, 0),
			),
		)).
		Param("urlprefix", u.urlPrefix).
		Path("ping").
		Fetch(ctx)
	if err != nil {
		return err
	}
	u.logger.Infoln("Notified", u, "with urlprefix", u.urlPrefix)
	return nil
}

func (u *Notifier) String() string {
	return fmt.Sprintf("Overcast (%s)", u.host)
}
