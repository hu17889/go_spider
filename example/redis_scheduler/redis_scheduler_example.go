package main

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/extension/scheduler"
    "net/url"
    "regexp"
    "strings"
    //"github.com/hu17889/go_spider/core/scheduler"
    "github.com/hu17889/go_spider/core/spider"
    "os"
)

type MyProcessor struct {
}

func (this *MyProcessor) Process(p *page.Page) {
    if !p.IsSucc() {
        mlog.LogInst().LogError(p.Errormsg())
        return
    }

    u, err := url.Parse(p.GetRequest().GetUrl())
    if err != nil {
        mlog.LogInst().LogError(err.Error())
        return
    }
    if !strings.HasSuffix(u.Host, "jiexieyin.org") {
        return
    }

    var urls []string
    query := p.GetHtmlParser()

    query.Find("a").Each(func(i int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        reJavascript := regexp.MustCompile("^javascript\\:")
        reLocal := regexp.MustCompile("^\\#")
        reMailto := regexp.MustCompile("^mailto\\:")
        if reJavascript.MatchString(href) || reLocal.MatchString(href) || reMailto.MatchString(href) {
            return
        }

        //处理相对路径
        var absHref string
        urlHref, err := url.Parse(href)
        if err != nil {
            mlog.LogInst().LogError(err.Error())
            return
        }
        if !urlHref.IsAbs() {
            urlPrefix := p.GetRequest().GetUrl()
            absHref = urlPrefix + href
            urls = append(urls, absHref)
        } else {
            urls = append(urls, href)
        }

    })

    p.AddTargetRequests(urls, "html")

}

func (this *MyProcessor) Finish() {

}

func main() {
    start_url := "http://www.jiexieyin.org"
    thread_num := uint(16)

    redisAddr := "127.0.0.1:6379"
    redisMaxConn := 10
    redisMaxIdle := 10

    proc := &MyProcessor{}

    sp := spider.NewSpider(proc, "redis_scheduler_example").
        //SetSleepTime("fixed", 6000, 6000).
        //SetScheduler(scheduler.NewQueueScheduler(true)).
        SetScheduler(scheduler.NewRedisScheduler(redisAddr, redisMaxConn, redisMaxIdle, true)).
        AddPipeline(pipeline.NewPipelineConsole()).
        SetThreadnum(thread_num)

    init := false
    for _, arg := range os.Args {
        if arg == "--init" {
            init = true
            break
        }
    }
    if init {
        sp.AddUrl(start_url, "html")
        mlog.LogInst().LogInfo("重新开始爬")
    } else {
        mlog.LogInst().LogInfo("继续爬")
    }
    sp.Run()
}
