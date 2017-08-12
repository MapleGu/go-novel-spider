package main

// NovelInfo 小说详情
type NovelInfo struct {
	CoverURL     string `json:"coverURL"`
	Name         string `json:"name"`
	Author       string `json:"author"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	Introduction string `json:"introduction"`
}

// ChapterChanStruct 章节
type ChapterChanStruct struct {
	Novel *NovelInfo
	Name  string
	URL   string
	Index int
}

// ChaptersListChanStruct 章节列表
type ChaptersListChanStruct struct {
	Novel *NovelInfo
	URL   string
}
