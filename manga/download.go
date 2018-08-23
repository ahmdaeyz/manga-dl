package manga

import (
	"fmt"
	"github.com/gocolly/colly"
	zip2 "github.com/pierrre/archivefile/zip"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	pageUrl string
	PageNum int
}

func (ch Chapter) DownloadByChapter(savePath string, cbz bool) int {
	downloaded, pageCount := 0, 0
	page := make(chan Page)
	c := colly.NewCollector(
		colly.UserAgent("firefox/61.0.1"),
	)
	var pages []Page
	c.OnHTML("div.read_img a img#image", func(element *colly.HTMLElement) {
		pageUrl := element.Attr("src")
		pages = append(pages, Page{pageUrl: pageUrl, PageNum: pageCount})
		pageCount++
	})
	for i := 1; i <= ch.numOfPages; i++ {
		c.Visit(strings.Replace(ch.url, "1.html", "", 1) + strconv.Itoa(i) + ".html")
	}
	ch.name = strings.Replace(ch.name, " ", "_", -1)
	CreateDirIfNotExist(savePath + "/" + ch.name)

	for i := 0; i < 30; i++ {
		go downPage(ch.name, savePath, page)
	}
	for i := 0; i < len(pages); i++ {
		page <- pages[i]
	}
	close(page)
	if cbz {
		defer os.RemoveAll(savePath + "/" + ch.name + "/")
		cbzPath, _ := os.Create(savePath + "/" + ch.name + ".cbz")
		zip2.Archive(savePath+"/"+ch.name, cbzPath, nil)
	}
	fmt.Println(ch.name, "downloaded.")
	time.Sleep(5 * time.Second)
	return downloaded
}

func downPage(dirName string, savePath string, chpages <-chan Page) bool {
	for page := range chpages {
		r, err := http.Get(page.pageUrl)
		defer r.Body.Close()
		f, ferr := os.Create(savePath + "/" + dirName + "/" + strconv.Itoa(page.PageNum+1) + ".jpg")
		_, imgErr := io.Copy(f, r.Body)
		if imgErr != nil {
			log.Fatal(imgErr)
		}
		f.Close()
		if err != nil && ferr != nil && imgErr != nil {
			return true
		}
	}
	return false
}

func (vol *Volume) DownloadByVolume(savePath string, cbz bool) int {
	numOfChDownloaded := 0
	volTitle := vol.VolNum
	if vol.VolNum == "Volume Not Available" {
		volTitle = "Volume 0"
	}
	pathToVolume := savePath + "/" + volTitle
	CreateDirIfNotExist(pathToVolume)
	chapters := vol.chapters
	for i := range chapters {
		numOfChDownloaded += chapters[i].DownloadByChapter(pathToVolume, false)
	}
	if cbz {
		defer os.RemoveAll(pathToVolume)
		cbzPath, _ := os.Create(pathToVolume + ".cbz")
		zip2.Archive(pathToVolume, cbzPath, nil)
	}
	fmt.Println(vol.VolNum, "downloaded.")
	return numOfChDownloaded
}
