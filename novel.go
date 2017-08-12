package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// NovelEncoding nove's html encode
const NovelEncoding = "gbk"
const website = "http://m.ybdu.com"

// Dir the place where novel save
const Dir = "/Users/guhuanian/Documents/novel"

// NovelChan 缓存 小说 url 详情
var NovelChan = make(chan string, 0)

// ChaptersChan 缓存 有 章节目录的 url
var ChaptersChan = make(chan *ChaptersListChanStruct, 1)

// ChapterURLsChan 缓存 待爬取得章节
var ChapterURLsChan = make(chan *ChapterChanStruct, 1)

// GetNovelURLFromDoc func find novel url
func GetNovelURLFromDoc(doc *goquery.Document) string {
	doc.Find(".line").Each(func(_ int, arg *goquery.Selection) {
		href := arg.Find("a").Eq(1).AttrOr("href", "")
		if href != "" {
			NovelChan <- website + href
		}
	})
	d := doc.Find(".page a").Eq(0)
	content := d.Text()
	if content != "下页" {
		return ""
	}
	return d.AttrOr("href", "")
}

// GetNovelInfoFromDoc func return novel info
func GetNovelInfoFromDoc(doc *goquery.Document) NovelInfo {
	coverURL := doc.Find(".block .block_img2 img").AttrOr("src", "")
	// var name, author, category, status, introduction string
	name := doc.Find(".block .block_txt2 h1 a").Text()

	d := doc.Find(".block .block_txt2 p")

	author := d.Eq(2).Text()
	author = strings.TrimPrefix(author, "作者：")

	category := d.Eq(3).Find("a").Text()
	category = strings.TrimPrefix(category, "全本")
	category = strings.TrimSuffix(category, "小说")

	status := d.Eq(4).Text()
	status = strings.TrimPrefix(status, "状态：")

	introduction := doc.Find(".intro_info").Text()

	return NovelInfo{
		CoverURL:     coverURL,
		Name:         name,
		Author:       author,
		Category:     category,
		Status:       status,
		Introduction: introduction,
	}
}
