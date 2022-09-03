package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {
	nr := flag.Int("n", 0, "Get the XKCD comic with this number")
	random := flag.Bool("r", false, "Get a random XKCD comic")
	updateLocal := flag.Bool("u", false, "Update locally stored XKCD comic descriptions")
	searchTerm := flag.String("s", "", "Searches all locally available XKCD comics for the provided term")
	flag.Parse()
	var comics []Info
	switch {
	case *nr > 0:
		comics = append(comics, getNr(*nr))
	case *random:
		comics = append(comics, findRandom())
	case *updateLocal:
		updateLocalFile()
	case *searchTerm != "":
		comics = append(comics, search(*searchTerm)...)
	}

	if len(comics) > 0 {
		fmt.Printf("%-5s\t%-40s\n", "Nr", "Title & Alternative Text")
		for _, comic := range comics {
			fmt.Printf("%-5d\t%-40s\n\t%s\n\t%s\n\n", comic.Num, comic.Title, comic.Alt, comic.Url)
		}
	}
}

func search(term string) []Info {
	var info []Info
	all, err := getAllFromLocal()
	if err != nil {
		log.Fatalf("Unable to retrieve comics from local storage: %s", err)
	}

	for _, comic := range all.Comics {
		s := strings.ToLower(comic.Title + comic.Alt + comic.Transcript)
		if strings.Contains(s, strings.ToLower(term)) {
			info = append(info, comic)
		}
	}
	return info
}

func findRandom() Info {
	rand.Seed(time.Now().UnixNano())
	latest := getLatestNr().Num
	return getNr(rand.Intn(latest-FirstComic) + FirstComic)
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
