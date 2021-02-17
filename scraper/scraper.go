package scraper

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"lightshot-scraper/util"
	"log"
	"net/http"
	"path"
	"sync"
)

const (
	userAgentKey      = "user-agent"
	userAgentValue    = "PostmanRuntime/7.26.8"
	imageIdSelector   = "#screenshot-image"
	imageUrlAttribute = "src"
	imagePermissions  = 0644
	symbolsCount      = 5
)

type lightshotScraper struct {
	baseUrl      string
	dirName      string
	client       http.Client
	urlGenerator util.UrlGenerator
}

type LightshotScraper interface {
	Scrap(context.Context, int) error
}

func New(
	dirName string,
	baseUrl string,
) LightshotScraper {
	return &lightshotScraper{
		dirName:      dirName,
		baseUrl:      baseUrl,
		client:       http.Client{},
		urlGenerator: util.NewUrlGenerator(),
	}
}

func (ls *lightshotScraper) Scrap(ctx context.Context, rounds int) error {
	util.CreateDirectoryIfNotExists(ls.dirName)
	var wg sync.WaitGroup

	for i := 0; i < rounds; i++ {
		wg.Add(1)
		imageHash := ls.urlGenerator.Generate(symbolsCount)
		imgSrc, err := ls.getImageSrc(ctx, imageHash)
		if err != nil {
			logrus.Errorf(
				"Error getting image src from page %v/%v: %v",
				ls.baseUrl,
				imageHash,
				err,
			)
			continue
		}
		go func() {
			defer wg.Done()
			if err := ls.saveImage(ctx, imgSrc); err != nil {
				logrus.Errorf("Error saving image %v: %v", imgSrc, err)
			}
		}()
	}

	wg.Wait()
}

func (ls *lightshotScraper) get(ctx context.Context, url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header[userAgentKey] = []string{userAgentValue}
	return ls.client.Do(req)
}

func (ls *lightshotScraper) getImageSrc(
	ctx context.Context,
	imageHash string,
) (string, error) {
	fullPath := ls.baseUrl + "/" + imageHash
	log.Println(fullPath)

	resp, err := ls.get(ctx, fullPath)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var imgSrc string
	doc.Find(imageIdSelector).Each(func(i int, s *goquery.Selection) {
		imgSrc, _ = s.Attr(imageUrlAttribute)
	})
	return imgSrc, nil
}

func (ls *lightshotScraper) saveImage(ctx context.Context, url string) error {
	errorCh := make(chan error, 1)
	defer close(errorCh)

	imageResp, err := ls.get(ctx, url)
	if err != nil {
		return err
	}

	image, err := ioutil.ReadAll(imageResp.Body)
	if err != nil {
		return err
	}

	fileName := path.Join(ls.dirName, path.Base(url))
	if err := ioutil.WriteFile(fileName, image, imagePermissions); err != nil {
		return err
	}
	return nil
}

func waitForCompletion(wg *sync.WaitGroup, errCh <-chan error) error {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case err := <-errCh:
		return err
	}
}
