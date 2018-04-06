package archive

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/datatogether/warc"
	"github.com/qri-io/cafs"
)

func PackagePathName(rawurl string) string {
	u, _ := url.Parse(rawurl)
	u.Scheme = ""

	// if u.Path == "" || u.Path == "/" || filepath.Base(u.Path) == "" {
	// if u.Path == "" || u.Path == "/" || filepath.Base(u.Path) == "" {
	if strings.HasSuffix(u.Path, "/") {
		// fmt.Println("adusting path:", u.Path, "->", u.Path+"index.html")
		u.Path += "index.html"
	} else if filepath.Ext(u.Path) == "" {
		// fmt.Println("adusting path:", u.Path, "->", u.Path+"/index.html")
		u.Path += "/index.html"
	}
	return strings.TrimPrefix(u.String(), "/")
}

func PackageRecords(urls []string, records warc.Records) (*cafs.Memdir, error) {
	pkg := cafs.NewMemdir("/")
	// cheap hack for now to only add files once, this should happen *much* earlier
	// in the archival process
	added := map[string]bool{}
	resRecs := records.FilterTypes(warc.RecordTypeResponse, warc.RecordTypeResource)
	for _, rec := range resRecs {
		body, err := rec.Body()
		if err != nil {
			fmt.Println("error getting body bytes:", err.Error())
			continue
		}

		// path := rw.Urlrw.RewriteString(rec.TargetURI())
		path := PackagePathName(rec.TargetURI())
		if added[path] {
			continue
		}
		fmt.Println(path)

		added[path] = true
		fmt.Println(path)
		pkg.AddChildren(cafs.NewMemfileBytes(path, body))
	}

	buf := &bytes.Buffer{}
	warc.WriteRecords(buf, records)

	indexBuf := &bytes.Buffer{}

	// TODO - improve cdxj index to add all fields required by cdxj spec
	// cdxi := make(cdxj.Index, len(urls))
	// for i, u := range urls {
	// 	// for now we're faking the actual index records
	// 	ir, err := cdxj.CreateRecord(&warc.Record{
	// 		Type: warc.RecordTypeRequest,
	// 		Headers: warc.Header{
	// 			warc.FieldNameWARCTargetURI: u,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	cdxi[i] = ir
	// }

	if err := RenderIndexTemplate(indexBuf, urls, records); err != nil {
		return nil, err
	}

	pkg.AddChildren(
		cafs.NewMemfileBytes("index.html", indexBuf.Bytes()),
		cafs.NewMemfileBytes("archive.warc", buf.Bytes()),
	)

	return pkg, nil
}
