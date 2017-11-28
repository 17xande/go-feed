package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	txtTemplate "text/template"
	"time"
)

type config struct {
	Title       string
	Link        string
	Author      string
	Type        string
	ImagePath   string
	Description string
	ItemsPath   string
	OwnerName   string
	OwnerEmail  string
}

type item struct {
	Title       string
	Enclosure   template.URL
	Description string
	Episode     int
	PubDate     string
	Duration    time.Duration
	Length      int64
}

func main() {
	var conf config
	addr := flag.String("addr", ":3000", "The address of the application")
	debug := flag.Bool("debug", false, "get more logs")
	ip := flag.String("path", "", "The path where the files to be hosted are")
	flag.Parse()

	conf.ItemsPath = *ip
	if len(conf.ItemsPath) > 1 && conf.ItemsPath[:2] == "~/" {
		u, err := user.Current()
		if err != nil {
			log.Println("error trying to get system user: ", err)
		}
		conf.ItemsPath = filepath.Join(u.HomeDir, conf.ItemsPath[2:])
	}

	if *debug {
		log.Println("debug mode enabled")
	}

	data, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		log.Fatal("error reading config file: ", err)
	}

	if *debug {
		log.Println("config file read")
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("error unmarshalling config file: ", err)
	}

	if *debug {
		log.Println("config file unmarshalled")
	}

	// read items path to list available items.
	if len(conf.ItemsPath) == 0 {
		log.Fatal("no item path provided. exiting program")
	}
	files, err := ioutil.ReadDir(conf.ItemsPath)
	if err != nil {
		log.Fatal("error trying to read items path of : ", conf.ItemsPath, err)
	}

	var items []item
	// filter out all files except for .mp3 and create items
	for _, file := range files {
		e := strings.Replace(file.Name(), " ", "%20", -1)
		if path.Ext(file.Name()) == ".mp3" {
			i := item{
				Title:       file.Name(),
				Enclosure:   template.URL("/podcasts/" + e),
				Description: "An item.",
				PubDate:     file.ModTime().Format(time.RFC1123Z),
				Length:      file.Size(),
			}
			items = append(items, i)
		}
	}

	if *debug {
		log.Println(len(items), " items read from directory")
	}

	http.HandleFunc("/", handlerHome)
	http.Handle("/podcasts/", http.StripPrefix("/podcasts/", http.FileServer(http.Dir(conf.ItemsPath))))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	http.HandleFunc("/podcast", func(w http.ResponseWriter, r *http.Request) {
		templ := txtTemplate.Must(txtTemplate.ParseFiles("web/templates/podcast.rss"))
		data := map[string]interface{}{
			"config": &conf,
			"items":  items,
			"host":   r.Host,
		}
		err := templ.Execute(w, data)
		if err != nil {
			log.Println("templ.Execute: ", err)
		}
	})

	log.Printf("Listening on port %s\n", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("could not start web server: ", err)
	}
}

func handlerTemplate(filename string, w io.Writer, data map[string]interface{}) {
	templ := template.Must(template.ParseFiles(filepath.Join("web/templates", filename)))
	err := templ.Execute(w, data)
	if err != nil {
		log.Println("templ.Execute: ", err)
	}
}

// Handles requests to the index page as well as any other requests
// that don't match any other paths
func handlerHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		log.Printf("Not found: %s", r.URL)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		log.Printf("Method not allowed: %s", r.URL)
		return
	}

	handlerTemplate("index.html", w, nil)
}
