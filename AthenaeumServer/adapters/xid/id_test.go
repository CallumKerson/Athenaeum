package xid_test

import (
	"bytes"
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/xid"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

func TestXidGeneratesProperly(t *testing.T) {
	buf := new(bytes.Buffer)
	adapter := xid.NewAdapter(logging.New(buf, "debug", "text"))

	//when
	id := adapter.Get()

	//then
	if len(id) != 20 {
		t.Errorf("Xid was the wrong length, expected length of %d but xid was %s", 20, id)
	}
}
