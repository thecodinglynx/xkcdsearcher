package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const FirstComic = 1
const LastComic = 2661 // as of August 21st
const XkcdUrl = "https://xkcd.com/%d/info.0.json"

type Info struct {
	Num              int
	Title            string
	Alt              string
	Year, Month, Day string
}

func getXkcd(nr int) Info {
	resp, err := http.Get(fmt.Sprintf(XkcdUrl, nr))
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
	return info
}
