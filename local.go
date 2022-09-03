package main

import (
	"encoding/json"
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
