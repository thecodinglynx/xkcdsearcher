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

// first attempts to retrieve the comic info from local storage, otherwise
// retrieves it from web
func getNr(nr int) Info {
	if local, err := getFromLocal(nr); err == nil && local != (Info{}) {
		fmt.Printf("Returning %d from local storage\n", nr)
		return local
	}
	fmt.Printf("Returning %d from web\n", nr)
	return getXkcd(fmt.Sprintf(XkcdUrl, nr))
}

// find and return locally stored XKCD comics that contain the provided search term
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
