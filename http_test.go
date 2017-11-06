package archive

import (
	// "bytes"
	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datatogether/warc"
)

func TestDoRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello!")
	}))

	req, _ := http.NewRequest("GET", s.URL, nil)
	reqr, resr, err := DoRequest(req, warc.Records{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	if reqr.Type != warc.RecordTypeRequest {
		t.Errorf("expected record to be request type")
	}

	if resr.Type != warc.RecordTypeResponse {
		t.Errorf("expected record to be response type")
	}
}
