package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const XkcdUrl = "https://xkcd.com/571/info.0.json"

type Info struct {
	Title string
	Alt   string
}

func getXkcd() {
	resp, err := http.Get(XkcdUrl)
	if err != nil {
		log.Fatalf("Unable to retrieve XKCD description from %s: %s", XkcdUrl, err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK status code returned: %v", resp.StatusCode)
	}
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("Unable to parse response to Info struct: %s", err)
	}
	resp.Body.Close()
	fmt.Printf("%-20s %s\n", info.Title, info.Alt)
}
