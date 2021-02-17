package util

import "flag"

type AppParams struct {
	OutputDirectoryPath  string
	PicturesToScrapCount int
}

func ParseArgs() AppParams {
	outputFlag := flag.String("o", "pics", "Directory for downloaded pictures storage")
	roundsCountFlag := flag.Int("n", 10, "Number of pictures to scrap")
	flag.Parse()

	return AppParams{
		OutputDirectoryPath:  *outputFlag,
		PicturesToScrapCount: *roundsCountFlag,
	}
}
