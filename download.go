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

func getAll() AllComics {
	var allComics AllComics
	jsonFile, err := os.Open(localFile)
	if err != nil {
		fmt.Printf("Unable to read local file %s - loading comics from web\n", localFile)
		allComics = getAllFromWeb()
		data, err := json.Marshal(allComics)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		if ioutil.WriteFile(localFile, data, os.ModePerm) != nil {
			log.Fatalf("Unable to write to file %s: %s", localFile, err)
		}
		return allComics
	}
	fmt.Printf("Reading comics from file %s\n", localFile)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &allComics)
	jsonFile.Close()

	return allComics
}

func getAllFromWeb() AllComics {
	var all AllComics
	for i := FirstComic; i < 5; i++ {
		all.Comics = append(all.Comics, getNr(i))
	}
	return all
}

func getNr(nr int) Info {
	return getXkcd(fmt.Sprintf(XkcdUrl, nr))
}

func getLatestNr() Info {
	return getXkcd(XkcdUrlLatest)
}

func getXkcd(url string) Info {
	resp, err := http.Get(url)
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
