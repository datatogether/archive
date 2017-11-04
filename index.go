package archive

import (
	"bytes"
	"html/template"
	"io"
	"time"

	"github.com/datatogether/cdxj"
)

func RenderIndexTemplate(w io.Writer, index cdxj.Index) error {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(w, "index.html", NewTemplateData("Data Together Archive", index))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf)
	return err
}

type TemplateData struct {
	Title     string
	Datestamp string
	Records   []*Record
}

type Record struct {
	Path      string
	Url       string
	Datestamp time.Time
}

func NewTemplateData(title string, index cdxj.Index) *TemplateData {
	recs := make([]*Record, len(index))
	for i, rec := range index {
		u, err := cdxj.UnSURTUrl(rec.Uri)
		if err != nil {
			continue
		}
		// p, err := cdxj.UnSURTPath(rec.Uri)
		// if err != nil {
		// 	continue
		// }

		recs[i] = &Record{
			Path: ArchivePathName(u),
			Url:  u,
		}
	}

	return &TemplateData{
		Title:     title,
		Records:   recs,
		Datestamp: time.Now().Format(time.UnixDate),
	}
}
