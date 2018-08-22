package manga

import (
	"github.com/gocolly/colly"
	"strings"
	"strconv"
	"os"
	"net/http"
	"io"
	"log"
	zip2 "github.com/pierrre/archivefile/zip"
	"fmt"
)
type Page struct{
	pageUrl string
	PageNum int
}

func (ch Chapter) DownloadByChapter(savePath string,cbz bool) int{
	downloaded,pageCount:=0,0
	c:=colly.NewCollector(
		colly.UserAgent("firefox/61.0.1"),
	)
	var pages []Page
	c.OnHTML("div.read_img a img#image", func(element *colly.HTMLElement) {
		pageUrl:=element.Attr("src")
		pages= append(pages,Page{pageUrl:pageUrl,PageNum:pageCount})
		pageCount++
	})
	for i:=1;i<=ch.numOfPages;i++ {
		c.Visit(strings.Replace(ch.url,"1.html","",1)+strconv.Itoa(i)+".html")
	}
	ch.name=strings.Replace(ch.name," ","_",-1)
	CreateDirIfNotExist(savePath+"/"+ch.name)
	for i:= range pages{
		pages[i].downPage(ch.name,savePath)
		downloaded++
	}
	if cbz{
		defer os.RemoveAll(savePath+"/"+ch.name+"/")
		cbzPath,_:=os.Create(savePath+"/"+ch.name+".cbz")
		zip2.Archive(savePath+"/"+ch.name,cbzPath,nil)
	}
	fmt.Println(ch.name,"downloaded.")
	return downloaded
}

func (page Page) downPage(dirName string,savePath string) bool{
	r,err:=http.Get(page.pageUrl)
	defer r.Body.Close()
	f,ferr:=os.Create(savePath+"/"+dirName+"/"+strconv.Itoa(page.PageNum)+".jpg")
	_,imgErr := io.Copy(f,r.Body)
	if imgErr!=nil{
		log.Fatal(imgErr)
	}
	f.Close()
	if err !=nil && ferr!=nil && imgErr!=nil{
		return true
	}
	return false
}

func (vol *Volume) DownloadByVolume(savePath string,cbz bool) int {
	numOfChDownloaded:=0
	pathToVolume:=savePath+"/"+"Volume_"+strconv.FormatFloat(vol.VolNum,'f',1,64)
	CreateDirIfNotExist(pathToVolume)
	chapters:=vol.chapters
	for i:= range chapters{
		numOfChDownloaded+=chapters[i].DownloadByChapter(pathToVolume,false)
	}
	if cbz{
		defer os.RemoveAll(pathToVolume)
		cbzPath,_:=os.Create(pathToVolume+".cbz")
		zip2.Archive(pathToVolume,cbzPath,nil)
	}
	fmt.Println("Volume",vol.VolNum,"downloaded.")
	return numOfChDownloaded
}


