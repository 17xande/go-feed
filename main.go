package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "The address of the application")
	debug := flag.Bool("debug", false, "get more logs")

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

	log.Printf("Listening on port %s\n", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("could not start web server: ", err)
	}
}
