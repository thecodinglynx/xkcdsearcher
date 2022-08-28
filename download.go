package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const FirstComic = 1
const XkcdUrl = "https://xkcd.com/%d/info.0.json"
const XkcdUrlLatest = "https://xkcd.com/info.0.json"
const localFile = "xkcd_info.json"

type AllComics struct {
	Comics []Info
}

type Info struct {
	Num              int
	Title            string
	Alt              string
	Year, Month, Day string
}

func getAllFromLocal() (AllComics, error) {
	var allComics AllComics
	jsonFile, err := os.Open(localFile)
	if err != nil {
		return allComics, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &allComics)

	return allComics, nil
}

func getFromLocal(nr int) (Info, error) {
	var info Info
	allComics, err := getAllFromLocal()
	if err != nil {
		return info, err
	}
	if c := allComics.Comics; len(c) > nr && c[nr-1] != (Info{}) {
		return c[nr-1], nil
	}
	return info, nil
}

func getFromWeb(nrs []int) AllComics {
	var all AllComics
	for i := range nrs {
		if info := getNr(nrs[i]); info != (Info{}) {
			all.Comics = append(all.Comics, info)
		}
	}
	return all
}

func writeToLocal(comics AllComics) {
	data, err := json.MarshalIndent(comics, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	if ioutil.WriteFile(localFile, data, os.ModePerm) != nil {
		log.Fatalf("Unable to write to file %s: %s", localFile, err)
	}
}

// first attempts to retrieve the comic info from local storage, otherwise
// retrieves it from web
func getNr(nr int) Info {
	if local, err := getFromLocal(nr); err == nil && local != (Info{}) {
		fmt.Printf("Returning %d from local storage\n\n", nr)
		return local
	}
	fmt.Printf("Returning %d from web\n\n", nr)
	return getXkcd(fmt.Sprintf(XkcdUrl, nr))
}

func getLatestNr() Info {
	return getXkcd(XkcdUrlLatest)
}

// TODO: make this run concurrently
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
