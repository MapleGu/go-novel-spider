package main

type NovelInfo struct {
	CoverURL     string `json:"coverURL"`
	Name         string `json:"name"`
	Author       string `json:"author"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	Introduction string `json:"introduction"`
}

type ChapterChanStruct struct {
	Novel *NovelInfo
	Name  string
	URL   string
	Index int
}

type ChaptersListChanStruct struct {
	Novel *NovelInfo
	URL   string
}
