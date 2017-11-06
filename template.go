package archive

import (
	"bytes"
	"html/template"
	"io"
	"time"

	"github.com/datatogether/warc"
)

func RenderIndexTemplate(w io.Writer, urls []string, records warc.Records) error {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(w, "index.html", NewTemplateData("Data Together Archive", urls, records))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf)
	return err
}

type TemplateData struct {
	Title     string
	Datestamp string
	Records   []*TemplateRecord
}

type TemplateRecord struct {
	Path      string
	Url       string
	Datestamp time.Time
}

func NewTemplateData(title string, urls []string, records warc.Records) *TemplateData {
	recs := make([]*TemplateRecord, len(urls))
	for i, u := range urls {
		recs[i] = &TemplateRecord{
			Path: "." + PackagePathName(u),
			Url:  u,
		}
	}

	return &TemplateData{
		Title:     title,
		Records:   recs,
		Datestamp: time.Now().Format(time.UnixDate),
	}
}
