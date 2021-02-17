package main

import (
	"context"
	"lightshot-scraper/scraper"
	"lightshot-scraper/util"
	"log"
)

func main() {
	params := util.ParseArgs()
	scrap := scraper.New(params.OutputDirectoryPath, "https://prnt.sc")
	if err := scrap.Scrap(context.Background(), params.PicturesToScrapCount); err != nil {
		log.Fatal(err)
	}
}
