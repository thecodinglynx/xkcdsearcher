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
	flag.Parse()
	var comics []Info
	if *nr > 0 {
		comics = append(comics, getNr(*nr))
	}
	if *random {
		comics = append(comics, findRandom())
	}
	fmt.Printf("%-5s %-40s %s\n", "Nr", "Title", "Alternative Text")
	for _, comic := range comics {
		fmt.Printf("%-5d %-40s %s\n", comic.Num, comic.Title, comic.Alt)
	}
}

func findRandom() Info {
	rand.Seed(time.Now().UnixNano())
	latest := getLatestNr().Num
	return getNr(rand.Intn(latest-FirstComic) + FirstComic)
}
