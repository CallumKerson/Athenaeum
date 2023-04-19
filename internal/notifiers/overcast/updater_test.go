package overcast

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	testURLPrefix := "https://athenaeum.test"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf("/ping?urlprefix=%s", url.QueryEscape(testURLPrefix)), req.URL.String())
		_, _ = rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	testNotifier := New(testURLPrefix, WithHost(server.URL))

	err := testNotifier.Notify(context.TODO())
	assert.NoError(t, err)
}
