package main

import (
	"html"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		l := strings.Replace(file.Name(), " ", "%20", -1)

		i := item{
			Title:       html.EscapeString(mp3File.Title()),
			Link:        template.URL("/podcasts/" + l),
			Description: html.EscapeString(d),
			PubDate:     file.ModTime().Format(time.RFC1123Z),
			Length:      file.Size(),
		}
		items = append(items, i)
	}

	return
}
