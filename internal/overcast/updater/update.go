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
	logger.Infoln("Will update Overcast for URL prefix", urlPrefix)
	return &Updater{host: OvercastHost, urlPrefix: urlPrefix, logger: logger}
}

func (u *Updater) Update(ctx context.Context) error {
	body := ""
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
		ToString(&body).
		Fetch(ctx)
	if err != nil {
		return err
	}
	u.logger.Infoln("Updated Overcast with urlprefix", u.urlPrefix, "response", body)
	return nil
}

func (u *Updater) String() string {
	return fmt.Sprintf("Overcast (%s)", u.host)
}
