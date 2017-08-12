package main

func main() {
	var url = "http://m.ybdu.com/book1/0/1/"
	// 根据 NovelChan 获取小说详情, 章节列表添加到 ChapterURLsChan
	go GetNovelInfoFromChannel()
	// 根据 ChapterChanStruct 获取章节内容
	go GetNovelChapterInfoChannel()
	go getNovelInfoFromChannelHandleChapters()
	for {
		GetNovelURL(url)
	}
}
