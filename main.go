package main

import (
	"context"
	"lightshot-scraper/scraper"
	"lightshot-scraper/util"
)

const lightshotUrl = "https://prnt.sc"

func main() {
	params := util.ParseArgs()
	scrap := scraper.New(params.OutputDirectoryPath, lightshotUrl)
	scrap.Scrap(context.Background(), params.PicturesToScrapCount)
}
