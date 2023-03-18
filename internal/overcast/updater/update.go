package updater

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

type Updater struct {
	host      string
	urlPrefix string
	logger    loggerrific.Logger
}

func New(urlPrefix string, logger loggerrific.Logger) *Updater {
	upd := &Updater{host: OvercastHost, urlPrefix: urlPrefix, logger: logger}
	logger.Infoln("Will update", upd, "for URL prefix", urlPrefix)
	return upd
}

func (u *Updater) Update(ctx context.Context) error {
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
	u.logger.Infoln("Updated", u, "with urlprefix", u.urlPrefix)
	return nil
}

func (u *Updater) String() string {
	return fmt.Sprintf("Overcast (%s)", u.host)
}
