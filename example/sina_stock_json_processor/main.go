// The example gets stock newses from site sina.com (http://live.sina.com.cn/zt/f/v/finance/globalnews1).
// The spider is continuous service.
// The stock api returns json result.
// It fetchs news at regular intervals that has been set in the config file.
// The result is saved in a file by PipelineFile.
package main

import (
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/spider"
    "fmt"
    "log"
    "strconv"
)

type MyPageProcesser struct {
    startNewsId int
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package simplejson (https://github.com/bitly/go-simplejson) is used to parse data of json.
func (this *MyPageProcesser) Process(p *page.Page) {
    if !p.IsSucc() {
        println(p.Errormsg())
        return
    }

    query := p.GetJson()
    status, err := query.GetPath("result", "status", "code").Int()
    if status != 0 || err != nil {
        log.Panicf("page is crawled error : errorinfo=%s : status=%d : startNewsId=%d", err.Error(), status, this.startNewsId)
    }
    num, err := query.GetPath("result", "pageStr", "pageSize").Int()
    if num == 0 || err != nil {
        // Add url of next crawl
        startIdstr := strconv.Itoa(this.startNewsId)
        p.AddTargetRequest("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id="+startIdstr+"&pagesize=10&dire=f", "json")
        return
    }

    var idint, nextid int
    var nextidstr string
    query = query.Get("result").Get("data")
    for i := 0; i < num; i++ {
        id, err := query.GetIndex(i).Get("id").String()
        if id == "" || err != nil {
            continue
        }
        idint, err = strconv.Atoi(id)
        if err != nil {
            continue
        }
        if idint <= this.startNewsId {
            break
        }
        if i == 0 {
            nextid = idint
            nextidstr = id
        }
        content, err := query.GetIndex(i).Get("content").String()
        if content == "" || err != nil {
            continue
        }
        time, err := query.GetIndex(i).Get("created_at").String()
        if err != nil {
            continue
        }

        p.AddField(id+"_id", id)
        p.AddField(id+"_content", content)
        p.AddField(id+"_time", time)
    }
    // Add url of next crawl
    this.startNewsId = nextid
    p.AddTargetRequest("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id="+nextidstr+"&pagesize=10&dire=f", "json")
    //println(p.GetTargetRequests())

}

func (this *MyPageProcesser) Finish() {
    fmt.Printf("TODO:before end spider \r\n")
}

func main() {
    // spider input:
    //  PageProcesser ;
    //  task name used in Pipeline for record;
    spider.NewSpider(NewMyPageProcesser(), "sina_stock_news").
        AddUrl("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=63621&pagesize=10&dire=f", "json"). // start url, html is the responce type ("html" or "json" or "jsonp" or "text")
        AddPipeline(pipeline.NewPipelineConsole()).                                                                                   // Print result to std output
        AddPipeline(pipeline.NewPipelineFile("/tmp/sinafile")).                                                                       // Print result in file
        OpenFileLog("/tmp").                                                                                                          // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".
        SetSleepTime("rand", 1000, 3000).                                                                                             // Sleep time between 1s and 3s.
        Run()
    //AddPipeline(pipeline.NewPipelineFile("/tmp/tmpfile")). // print result in file
}
