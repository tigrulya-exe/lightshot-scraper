package util

import "flag"

type AppParams struct {
	OutputDirectoryPath  string
	PicturesToScrapCount int
}

func ParseArgs() AppParams {
	outputFlag := flag.String("o", "pics", "The path where images will be stored")
	roundsCountFlag := flag.Int("n", 10, "Number of pictures to scrape")
	flag.Parse()

	return AppParams{
		OutputDirectoryPath:  *outputFlag,
		PicturesToScrapCount: *roundsCountFlag,
	}
}
