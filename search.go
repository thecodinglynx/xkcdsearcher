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
		fmt.Printf("%-5s %-40s %s\n", "Nr", "Title", "Alternative Text")
		for _, comic := range comics {
			fmt.Printf("%-5d %-40s %s\n", comic.Num, comic.Title, comic.Alt)
		}
	}
}

func findRandom() Info {
	rand.Seed(time.Now().UnixNano())
	latest := getLatestNr().Num
	return getNr(rand.Intn(latest-FirstComic) + FirstComic)
}

func updateLocalFile() {
	latestComicId := getLatestNr().Num
	allComics, _ := getAllFromLocal()
	var localUpToDate bool
	var latestLocalComicId int
	for _, v := range allComics.Comics {
		if v.Num > latestLocalComicId {
			latestLocalComicId = v.Num
		}
		if v.Num == latestComicId {
			localUpToDate = true
			break
		}
	}
	if localUpToDate {
		fmt.Println("Local file already up to date, no need to update.")
		return
	}
	fmt.Printf("Local file needs updating.\nLatest available comic:\t%d\nLatest stored comic:\t%d", latestComicId, latestLocalComicId)
	// TODO: determine diff between web and local, get diff from web and update local file
}
