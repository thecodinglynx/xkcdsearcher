package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const FirstComic = 1
const XkcdUrl = "https://xkcd.com/%d/info.0.json"
const XkcdUrlLatest = "https://xkcd.com/info.0.json"

func getFromWeb(nrs []int) AllComics {
	var all AllComics
	for i := range nrs {
		if info := getNr(nrs[i]); info != (Info{}) {
			all.Comics = append(all.Comics, info)
		}
	}
	return all
}

func getLatestNr() Info {
	return getXkcd(XkcdUrlLatest)
}

// TODO : make this run concurrently
func getXkcd(url string) Info {
	var info Info
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to retrieve XKCD description from %s: %s", XkcdUrl, err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK status code returned for url %s: %v - skipping\n", url, resp.StatusCode)
		return info
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("Unable to parse response to Info struct: %s", err)
	}
	resp.Body.Close()
	return info
}
