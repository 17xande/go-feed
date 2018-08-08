package main

import (
	"bytes"
	"encoding/xml"
	"html"
	"html/template"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dhowden/tag"
)

const extFilter = ".mp3"

type item struct {
	Title       string
	Link        template.URL
	Description string
	Author      string
	Episode     int
	PubDate     string
	Duration    time.Duration
	Length      int64
}

func genItems(files []os.FileInfo, dir string) (items []item) {
	for _, file := range files {
		// Ignore non .mp3 files.
		if path.Ext(file.Name()) != extFilter {
			continue
		}

		f, err := os.Open(filepath.Join(dir, file.Name()))
		defer f.Close()
		if err != nil {
			log.Println("Error trying to open file:", file.Name(), err)
			continue
		}

		mp3File, err := tag.ReadFrom(f)
		if err != nil {
			log.Println("Error trying to read ID3 tag info:", file.Name(), err)
			continue
		}

		d := "Speaker: " + mp3File.Artist()
		l := url.PathEscape(file.Name())
		buf := new(bytes.Buffer)
		xml.EscapeText(buf, []byte(l))

		title := html.EscapeString(mp3File.Title())
		if title == "" {
			title = path.Base(file.Name())
		}

		i := item{
			Title:       title,
			Link:        template.URL("/podcasts/" + buf.String()),
			Description: html.EscapeString(d),
			PubDate:     file.ModTime().Format(time.RFC1123Z),
			Length:      file.Size(),
		}
		items = append(items, i)
	}

	return
}
