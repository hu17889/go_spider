//
package main

import (
    "fmt"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/spider"
    "strings"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
    query := p.GetHtmlParser()

    name := query.Find(".lemmaTitleH1").Text()
    name = strings.Trim(name, " \t\n")

    summary := query.Find(".card-summary-content .para").Text()
    summary = strings.Trim(summary, " \t\n")

    // the entity we want to save by Pipeline
    p.AddField("name", name)
    p.AddField("summary", summary)
}

func main() {
    // spider input:
    //  PageProcesser ;
    //  config path(default: WD/etc/main.conf);
    //  task name used in Pipeline for record;
    spider.NewSpider(NewMyPageProcesser(), "", "sina_stock_news").
        AddUrl("https://github.com/hu17889?tab=repositories", "html"). // start url, html is the responce type ("html" or "json")
        AddPipeline(pipeline.NewPipelineConsole()).                    // print result on screen
        SetThreadnum(3).                                               // crawl request by three Coroutines
        Run()
}
