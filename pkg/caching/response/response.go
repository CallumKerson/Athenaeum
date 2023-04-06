package response

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"time"
)

type Response struct {
	Value      []byte
	Header     http.Header
	Expiration time.Time
	LastAccess time.Time
}

func BytesToResponse(b []byte) Response {
	var r Response
	dec := gob.NewDecoder(bytes.NewReader(b))
	_ = dec.Decode(&r)

	return r
}

func (r *Response) Bytes() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	_ = enc.Encode(r)

	return b.Bytes()
}
