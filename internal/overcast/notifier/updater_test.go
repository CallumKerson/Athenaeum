package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"
)

func TestUpdater(t *testing.T) {
	testURLPrefix := "https://athenaeum.test"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf("/ping?urlprefix=%s", url.QueryEscape(testURLPrefix)), req.URL.String())
		_, _ = rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	testNotifier := &Notifier{host: server.URL, urlPrefix: testURLPrefix, logger: tlogger.NewTLogger(t)}

	err := testNotifier.Notify(context.TODO())
	assert.NoError(t, err)
}
