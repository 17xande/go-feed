package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type config struct {
	Title     string
	Link      string
	Author    string
	ImagePath string
}

func main() {
	addr := flag.String("addr", ":3000", "The address of the application")
	debug := flag.Bool("debug", false, "get more logs")
	var conf config

	if *debug {
		log.Println("debug mode enabled")
	}

	data, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		log.Fatal("error reading config file: ", err)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("error unmarshalling config file: ", err)
	}

	http.HandleFunc("/", handlerHome)

	log.Printf("Listening on port %s\n", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("could not start web server: ", err)
	}
}

func handlerTemplate(filename string, w io.Writer, data map[string]interface{}) {
	templ := template.Must(template.ParseFiles(filepath.Join("web/templates", filename)))
	templ.Execute(w, data)
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
