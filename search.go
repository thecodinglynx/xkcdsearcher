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
	all := flag.Bool("a", false, "Get all XKCD comics")
	flag.Parse()
	var comics []Info
	switch {
	case *nr > 0:
		comics = append(comics, getNr(*nr))
	case *random:
		comics = append(comics, findRandom())
	case *all:
		allComics := getAll()
		for _, v := range allComics.Comics {
			fmt.Printf("%-5d %-40s %s\n", v.Num, v.Title, v.Alt)
		}
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
