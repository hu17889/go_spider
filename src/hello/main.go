package main

import (
	"page_processer"
	"common/page"
	"spider"
)


type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}
	
// 解析page内容，提取a标签等,用来
func (this *MyPageProcesser) Process(p *page.Page) {
}


func main() {
	var pagepro page_processer.PageProcesser
	pagepro = NewMyPageProcesser()
	var sp *spider.Spider
	sp = spider.NewSpider(pagepro)

	sp.Get("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1", "json")

	urls := []string{
		"http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1",
		"http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1",
	}
	sp.GetAll(urls, "json")
}