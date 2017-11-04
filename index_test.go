package archive

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/datatogether/cdxj"
	"github.com/datatogether/warc"
)

func TestRenderIndexTemplate(t *testing.T) {
	idx := cdxj.Index{
		&cdxj.Record{Uri: "(com,cnn,)/world", Timestamp: time.Time{}, RecordType: warc.RecordTypeRequest},
	}
	buf := &bytes.Buffer{}
	RenderIndexTemplate(buf, idx)

	fmt.Println(buf.String())
}
