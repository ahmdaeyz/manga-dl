package manga

import (
	"github.com/gocolly/colly"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Manga struct {
	url     string
	volumes []Volume
}

type Volume struct {
	VolNum   string
	startCh  float64
	endCh    float64
	chapters []Chapter
}

type Chapter struct {
	name       string
	url        string
	ChapterNum float64
	numOfPages int
	volNum     string
}

func GetManga(url string) Manga {
	/*
	*1*gets volumes data : num,starting,ending ch,...
	*2*assignes every chapter to its volume
	*returns a manga type
	 */

	//*1*
	var volumesData []string
	var volumes []Volume
	c := colly.NewCollector(
		colly.UserAgent("firefox/61.0.1"))
	c.OnHTML("div.left div#chapters div.slide h3.volume", func(element *colly.HTMLElement) {
		volumeDatum := element.Text
		volumesData = append(volumesData, volumeDatum)
	})
	c.Visit(url)
	rand.Seed(time.Now().Unix())
	//volume data : smth similar to volume 24 30-40
	for _, str := range volumesData {
		var startCh, endCh float64
		index := strings.Index(str, "C")
		volNum := strings.Replace(str, str[index-1:], "", -1)
		volumeBounds := strings.Replace(str, str[:index], "", -1)
		volumeBounds = strings.Replace(volumeBounds, "Chapter ", "", -1)
		chaps := strings.Split(volumeBounds, " - ")
		startCh, _ = strconv.ParseFloat(chaps[0], 64)
		if len(chaps) == 2 {
			endCh, _ = strconv.ParseFloat(chaps[1], 64)
		}
		volumes = append(volumes, Volume{VolNum: volNum, startCh: startCh, endCh: endCh})
	}
	//*2*
	chapters := GetChapters(url)
	for i := range volumes {
		for j := range chapters {
			if chapters[j].ChapterNum >= volumes[i].startCh && chapters[j].ChapterNum <= volumes[i].endCh {
				chapters[j].setVolumeNum(volumes[i].VolNum)
				volumes[i].appendChToVolume(chapters[j])
			}
		}
	}
	manga := Manga{url: url, volumes: volumes}
	return manga
}

func GetChapters(url string) []Chapter {
	/*
	*gets chapter's name along with url and number of pages
	*there are two because some of the chapters are stored in the h3 tag and others in the h4 try
	*inspecting a manga page from the site like:"http://fanfox.net/manga/berserk/"
	 */
	var chapters []Chapter
	c := colly.NewCollector(
		colly.UserAgent("firefox/61.0.1"))
	c.OnHTML("div#chapters ul.chlist li div h3 a.tips", func(element *colly.HTMLElement) {
		rg, _ := regexp.Compile("[0-9]{1,4}(\\.[0-9])?")
		chapterNum, _ := strconv.ParseFloat(rg.FindString(element.Text), 64)
		chapters = append(chapters, Chapter{name: element.Text, url: "http:" + element.Attr("href"), ChapterNum: chapterNum, numOfPages: NumberOfPages("http:" + element.Attr("href"))})
	})
	c.OnHTML("div#chapters ul.chlist li div h4 a.tips", func(element *colly.HTMLElement) {
		rg, _ := regexp.Compile("[0-9]{1,4}(\\.[0-9])?")
		chapterNum, _ := strconv.ParseFloat(rg.FindString(element.Text), 64)
		chapters = append(chapters, Chapter{name: element.Text, url: "http:" + element.Attr("href"), ChapterNum: chapterNum, numOfPages: NumberOfPages("http:" + element.Attr("href"))})
	})
	c.Visit(url)

	return chapters
}
func NumberOfPages(url string) int {
	//anon function checks if the slice contains the scrapped value
	c := colly.NewCollector(
		colly.UserAgent("firefox/61.0.1"),
	)
	var pages []int
	c.OnHTML("div.r.m div.l select.m option", func(element *colly.HTMLElement) {
		page, _ := strconv.Atoi(element.Attr("value"))
		if !func(arr []int, ele int) bool {
			for _, element := range arr {
				if element == ele {
					return true
					break
				}
			}
			return false
		}(pages, page) {
			pages = append(pages, page)
		}
	})
	c.Visit(url)
	sort.Ints(pages)
	return pages[len(pages)-1]
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0707)
		if err != nil {
			panic(err)
		}
	}
}
func (ch *Chapter) setVolumeNum(volNum string) {
	ch.volNum = volNum
}
func (vol *Volume) appendChToVolume(ch Chapter) {
	vol.chapters = append(vol.chapters, ch)
}
func (manga *Manga) setUrl(url string) {
	manga.url = url
}
func (manga Manga) GetVolumes() []Volume {
	return manga.volumes
}
func (vol Volume) GetChapters() []Chapter {
	return vol.chapters
}
func (ch *Chapter) setNumOfPages() {
	ch.numOfPages = NumberOfPages(ch.url)
}
