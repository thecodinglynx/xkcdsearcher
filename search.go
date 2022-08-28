package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	nr := flag.Int("n", 0, "Get the XKCD comic with this number")
	random := flag.Bool("r", false, "Get a random XKCD comic")
	updateLocal := flag.Bool("u", false, "Update locally stored XKCD comic descriptions")
	flag.Parse()
	var comics []Info
	switch {
	case *nr > 0:
		comics = append(comics, getNr(*nr))
	case *random:
		comics = append(comics, findRandom())
	case *updateLocal:
		updateLocalFile()
	}
	if len(comics) > 0 {
		fmt.Printf("%-5s\t%-40s\n", "Nr", "Title & Alternative Text")
		for _, comic := range comics[:10] {
			fmt.Printf("%-5d\t%-40s\n\t%s\n\n", comic.Num, comic.Title, comic.Alt)
		}
	}
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
		fmt.Printf("Local file already up to date, no need to update.\n\n")
		return allLocal.Comics
	}
	fmt.Printf("Latest available comic:\t\t%d\n", latestComicId)
	fmt.Printf("Latest locally stored comic:\t%d\n", latestLocalComicId)
	fmt.Printf("Downloading %d missing comics from web\n\n", len(missing))
	allWeb := getFromWeb(missing)
	combinedComics := append(allLocal.Comics, allWeb.Comics...)
	writeToLocal(AllComics{combinedComics})
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
