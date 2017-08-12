package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/iconv.v1"
)

// 获取小说

// GetNovelURL get novel url from category
func GetNovelURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("function getNovelURL's url cannot be blank")
	}

	doc, err := GetHTMLContentFromURL(url)
	if err != nil {
		return "", err
	}
	// 写入 NovelChan
	// GetNovelURLFromDoc 来自 novel.go
	url = GetNovelURLFromDoc(doc)

	return url, nil
}

// GetNovelInfoFromChannel get novel info from channel
func GetNovelInfoFromChannel() {
	for {
		url := <-NovelChan
		doc, err := GetHTMLContentFromURL(url)
		if err != nil {
			log.Println("获取小说详情失败 -> ", url, err)
			sleep()
			NovelChan <- url
			return
		}
		info := getNovelInfoFromChannelHandleInfo(doc)

		href := doc.Find(".cover .ablum_read").First().Find("span").Eq(1).Find("a").AttrOr("href", "")
		ChaptersChan <- &ChaptersListChanStruct{
			Novel: info,
			URL:   href,
		}
	}
}

// GetNovelChapterInfoChannel 获取小说内容
func GetNovelChapterInfoChannel() {
	for {
		data := <-ChapterURLsChan
		doc, err := GetHTMLContentFromURL(data.URL)
		if err != nil {
			log.Println("获取小说章节内容失败 -> ", data, err)
			sleep()
			ChapterURLsChan <- data
		}

		chapterPath := fmt.Sprintf("%s/%s-%s-%s/chapters", Dir, data.Novel.Category, data.Novel.Name, data.Novel.Author)
		os.Mkdir(chapterPath, 0700)
		f, err := os.Create(chapterPath + "/" + strconv.Itoa(data.Index) + "_" + data.Name + ".txt")
		if err != nil {
			log.Println("创建小说章节文件失败 ->", err, data)
			sleep()
			return
		}
		defer f.Close()
		content := doc.Find(".content .txt").Text()
		f.WriteString(content)
	}
}

func getNovelInfoFromChannelHandleChapters() {
	for {
		data := <-ChaptersChan
		doc, err := GetHTMLContentFromURL(data.URL)
		if err != nil {
			log.Println("获取小说章节列表失败 -> ", data, err)
			sleep()
			ChaptersChan <- data
			return
		}

		doc.Find(".chapter li a").Each(func(i int, d *goquery.Selection) {
			url := d.AttrOr("href", "")
			if url == "" {
				return
			}
			url = data.URL + url
			name := d.Text()
			ChapterURLsChan <- &ChapterChanStruct{
				Novel: data.Novel,
				Name:  name,
				URL:   url,
				Index: i,
			}
		})
	}
}
func getNovelInfoFromChannelHandleInfo(doc *goquery.Document) *NovelInfo {

	info := GetNovelInfoFromDoc(doc)

	novelPath := fmt.Sprintf("%s/%s-%s-%s", Dir, info.Category, info.Name, info.Author)

	os.Mkdir(novelPath, 0700)
	downImage(info.CoverURL, novelPath+"/cover.jpg")
	infoPath := fmt.Sprintf("%s/%s", novelPath, "info.json")
	f, err := os.Create(infoPath)
	if err != nil {
		log.Println("创建小说详情 info.json 失败 ->", err)
		sleep()
	}
	defer f.Close()

	infoString, _ := json.Marshal(info)
	f.Write(infoString)

	return &info
}

// GetHTMLContentFromURL get content from url and return goquery.Document
func GetHTMLContentFromURL(url string) (*goquery.Document, error) {
	if url == "" {
		return nil, errors.New("function GetHTMLContentFromURL's url cannot be blank")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cd, err := iconv.Open("utf-8", NovelEncoding)
	if err != nil {
		return nil, err
	}
	defer cd.Close()

	r := iconv.NewReader(cd, resp.Body, 0)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func downImage(url, path string) {
	if url == "" {
		return
	}
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	image, err := os.Create(path)
	if err != nil {
		return
	}
	defer image.Close()
	image.Write(data)
	return
}

// 发生错误，暂停 n 秒
func sleep() {
	time.Sleep(3 * time.Second)
}
