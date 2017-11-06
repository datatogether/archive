package archive

import (
	"net/http"
	"net/url"

	"github.com/datatogether/resources"
	"github.com/datatogether/rewrite"
	"github.com/datatogether/warc"
)

func ArchiveUrl(req *http.Request, rw *rewrite.WarcRecordRewriter, records warc.Records) (addRecords warc.Records, err error) {
	var reqr, resr *warc.Record
	var resources warc.Records
	reqr, resr, err = DoRequest(req, records)
	if err != nil {
		return
	}
	addRecords = warc.Records{reqr, resr}

	resources, err = extractResourceRecords(resr, records)
	if err != nil {
		return
	}

	addRecords = append(addRecords, resources...)
	for _, rec := range addRecords {
		resr, err = rw.RewriteRecord(rec)
		if err != nil {
			return
		}
	}

	return
}

func extractResourceRecords(resr *warc.Record, records warc.Records) (warc.Records, error) {
	rrecs := warc.Records{}
	ext := resources.NewExtractor()
	urls, err := ext.ExtractResponseUrls(resr)
	if err != nil {
		return nil, err
	}

	reqUrl, err := url.Parse(resr.TargetUri())
	if err != nil {
		return nil, err
	}

	for _, u := range urls {
		abs, err := reqUrl.Parse(u)
		if err != nil {
			continue
		}

		// skip if we've already archived this url
		if records.TargetUriRecord(abs.String(), warc.RecordTypeResponse, warc.RecordTypeResource) != nil ||
			rrecs.TargetUriRecord(abs.String(), warc.RecordTypeResponse, warc.RecordTypeResource) != nil {
			continue
		}

		req, err := http.NewRequest("GET", abs.String(), nil)
		if err != nil {
			// todo - log warning
			continue
		}
		rq, rs, err := DoRequest(req, records)
		if err != nil {
			return rrecs, err
		}
		rrecs = append(rrecs, rq, rs)
	}

	return rrecs, nil
}
