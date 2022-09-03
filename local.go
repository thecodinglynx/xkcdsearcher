package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const localFile = "xkcd_info.json"

type AllComics struct {
	Comics []Info
}

type Info struct {
	Num              int
	Title            string
	Alt              string
	Transcript       string
	Url              string `json:"img"`
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
	for _, info := range allComics.Comics {
		if info.Num == nr {
			return info, nil
		}
	}
	return info, nil
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

func updateLocalFile() []Info {
	latestComicId := getLatestNr().Num
	allLocal, _ := getAllFromLocal()
	var latestLocalComicId int

	// create bool slice for all comics, values default to false
	haveLocal := make([]bool, latestComicId)

	// funny guy - XKCD 404 doesn't exist
	haveLocal[404-1] = true
	for _, v := range allLocal.Comics {
		haveLocal[v.Num-1] = true
		if v.Num > latestLocalComicId {
			latestLocalComicId = v.Num
		}
	}
	missing := findMissing(haveLocal)
	if len(missing) <= 0 {
		fmt.Printf("All currently available Comics (%d) cached locally already, no need to update.\n\n", latestComicId)
		return allLocal.Comics
	}
	fmt.Printf("Latest available comic:\t\t%d\n", latestComicId)
	fmt.Printf("Latest locally stored comic:\t%d\n", latestLocalComicId)
	allWeb := getFromWeb(missing)
	combinedComics := append(allLocal.Comics, allWeb.Comics...)
	writeToLocal(AllComics{combinedComics})
	fmt.Printf("Successfully downloaded %d missing comics from web and updated local cache\n\n", len(missing))
	return combinedComics
}

func findMissing(haveLocal []bool) []int {
	var missing []int
	for i, have := range haveLocal {
		if !have {
			missing = append(missing, i+1)
		}
	}
	return missing
}
