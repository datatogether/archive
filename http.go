package archive

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/datatogether/ffi"
	"github.com/datatogether/warc"
)

// DoRequest performs an archival HTTP request, adding request & response records to
// the given records pointer, returning the response record & a list of all dependant resources
func DoRequest(req *http.Request, records warc.Records) (reqr, resr *warc.Record, err error) {
	// don't perform requests for urls already in this list of archives
	if rec := records.TargetUriRecord(req.URL.String(), warc.RecordTypeResponse, warc.RecordTypeResource); rec != nil {
		return
	}

	reqr = RequestRecord(req)
	reqr.Headers[warc.FieldNameWARCDate] = time.Now().Format(time.RFC3339)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	resr, err = HttpResponseRecord(res)
	if err != nil {
		return
	}

	return
}

func RequestRecord(req *http.Request) *warc.Record {
	body := contentFromHttpRequest(req)
	return &warc.Record{
		Type: warc.RecordTypeRequest,
		Headers: map[string]string{
			warc.FieldNameContentType:   "application/http; msgtype=request",
			warc.FieldNameWARCRecordID:  warc.NewUuid(),
			warc.FieldNameWARCTargetURI: req.URL.String(),
		},
		Content: bytes.NewBuffer(body),
	}
}

func contentFromHttpRequest(req *http.Request) []byte {
	buf := &bytes.Buffer{}

	if err := warc.WriteRequestStatusAndHeaders(buf, req); err != nil {
		return buf.Bytes()
	}

	// buf.WriteString(fmt.Sprintf("%s / %s\r\n", req.Method, req.Proto))
	// buf.WriteString(fmt.Sprintf("Host: %s\r\n", req.Host))
	// if err := writeHttpHeaders(buf, req.Header); err != nil {
	// 	fmt.Println("error writing to buffer:", err.Error())
	// }

	// buf.WriteString(fmt.Sprintf("User-Agent: %s\r\n", req.UserAgent()))
	// TODO - finish

	return buf.Bytes()
}

// HttpResponseRecord creates a record from an HTTP response
func HttpResponseRecord(res *http.Response) (*warc.Record, error) {
	raw, sanitized, mimetype, err := SanitizeResponse(res)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	warc.WriteHttpHeaders(buf, res.Header)
	buf.WriteString("\r\n")
	buf.Write(sanitized)

	resr := &warc.Record{
		Type: warc.RecordTypeResponse,
		Headers: map[string]string{
			warc.FieldNameWARCPayloadDigest:         warc.Sha1Digest(raw),
			warc.FieldNameContentType:               "application/http; msgtype=response",
			warc.FieldNameWARCRecordID:              warc.NewUuid(),
			warc.FieldNameWARCIdentifiedPayloadType: mimetype,
			warc.FieldNameWARCTargetURI:             res.Request.URL.String(),
		},
		Content: buf,
	}
	return resr, nil
}

func SanitizeResponse(res *http.Response) (raw, sanitized []byte, mimetype string, err error) {
	defer res.Body.Close()
	raw, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	mimetype = ffi.DetectContentType(res.Request.URL.String(), raw)
	sanitized, err = warc.Sanitize(mimetype, raw)
	return
}
