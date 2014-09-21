go_spider
=========
[![Build Status](https://travis-ci.org/hu17889/go_spider.svg)](https://travis-ci.org/hu17889/go_spider)


**v1.0** 完成基本框架与完整爬虫功能。

## 简介


本项目基于golang开发，是一个开放的垂直领域的爬虫引擎，主要希望能将各个功能模块区分开，方便使用者重新实现子模块，进而构建自己垂直方方向的爬虫。

本项目将爬虫的各个功能流程区分成Spider模块（主控），Downloader模块（下载器），PageProcesser模块（页面分析），Scheduler模块（任务队列），Pipeline模块（结果输出）；


**执行过程简述**：

1. Spider从Scheduler中获取包含待抓取url的Request对象，启动一个协程，一个协程执行一次爬取过程，此处我们把协程也看成Spider，Spider把Request对象传入Downloader，Downloader下载该Request对象中url所对应的页面或者其他类型的数据，生成Page对象；
2. Spider调用PageProcesser模块解析页面中的数据，并存入Page对象中的PageItems中（以Key-Value对的形式保存），同时存入解析结果中的待抓取链接，Spider会将待抓取链接存入Scheduler模块中的Request队列中；
3. Spider调用Pipeline模块输出Page中的PageItems的结果;
4. 执行步骤1，直至所有链接被处理完成，则Spider被挂起等待下一个待抓取链接或者终止。


![image](https://github.com/hu17889/doc/blob/master/go_spider/img/project.png)


执行过程相应的Spider核心代码，代码代表一次爬取过程：

``` Go
// core processer
func (this *Spider) pageProcess(req *request.Request) {
    p := this.pDownloader.Download(req)
    if p == nil {
        return
    }

    this.pPageProcesser.Process(p)
    for _, req := range p.GetTargetRequests() {
        this.addRequest(req)
    }

    // output
    if !p.GetSkip() {
        for _, pip := range this.pPiplelines {
            pip.Process(p.GetPageItems(), this)
        }
    }

    this.sleep()
}
```


## 安装

```
go get github.com/hu17889/go_spider
go get github.com/PuerkitoBio/goquery
go get github.com/bitly/go-simplejson
```


## 简单示例

示例中在main包中实现了爬虫创建，初始化，以及PageProcesser模块的继承实现。
示例的功能是爬取[https://github.com/hu17889?tab=repositories](https://github.com/hu17889?tab=repositories)下面的项目以及项目详情页的相关信息，并将内容输出到标准输出。

``` Go
//
package main

/*
Packages must be imported:
    "core/common/page"
    "core/spider"
Pckages may be imported:
    "core/pipeline": scawler result persistent;
    "github.com/PuerkitoBio/goquery": html dom parser.
*/
import (
    "github.com/PuerkitoBio/goquery"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/spider"
    "strings"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
    query := p.GetHtmlParser()
    var urls []string
    query.Find("h3[class='repo-list-name'] a").Each(func(i int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        urls = append(urls, "http://github.com/"+href)
    })
    // these urls will be saved and crawed by other coroutines.
    p.AddTargetRequests(urls, "html")

    name := query.Find(".entry-title .author").Text()
    name = strings.Trim(name, " \t\n")
    repository := query.Find(".entry-title .js-current-repository").Text()
    repository = strings.Trim(repository, " \t\n")
    //readme, _ := query.Find("#readme").Html()
    if name == "" {
        p.SetSkip(true)
    }
    // the entity we want to save by Pipeline
    p.AddField("author", name)
    p.AddField("project", repository)
    //p.AddField("readme", readme)
}

func main() {
    // spider input:
    //  PageProcesser ;
    //  config path(default: WD/etc/main.conf);
    //  task name used in Pipeline for record;
    spider.NewSpider(NewMyPageProcesser(), "", "TaskName").
        AddUrl("https://github.com/hu17889?tab=repositories", "html"). // start url, html is the responce type ("html" or "json")
        AddPipeline(pipeline.NewPipelineConsole()).                    // print result on screen
        SetThreadnum(3).                                               // crawl request by three Coroutines
        Run()
}

```


## 模块

### [Spider](http://godoc.org/github.com/hu17889/go_spider/core/spider)

**功能**：用户一般无需自己实现。完成爬虫初始化，如加入各个默认子模块，管理并发，调度其他模块以及相关参数设置。


### [Downloader](http://godoc.org/github.com/hu17889/go_spider/core/downloader)

**功能**：用户一般无需自己实现。Spider从Scheduler的Request队列中获取包含待抓取url的Request对象，传入Downloader，Downloader下载该Request对象中的url所对应的页面或者其他类型的数据，现在支持html和json两种结果类型或者无结果类型，生成Page对象，同时找到下载结果所对应的解析go包并生成解析器存入Page对象中，如html是[goquery包](https://github.com/PuerkitoBio/goquery)，json数据是[simplejson包](https://github.com/bitly/go-simplejson/blob/master/simplejson.go)。


### [PageProcesser](http://godoc.org/github.com/hu17889/go_spider/core/page_processer)

**功能**：用户必须实现此模块。这个模块主要做页面解析，用户需要在此处获取有用数据和下一步爬取的链接。PageProcesser的前后实现步骤如下：Spider调用PageProcesser模块解析页面中的数据，并存入Page对象中的PageItems对象中（以Key-Value对的形式保存），同时存入解析结果中的待抓取链接，Spider会将待抓取链接存入Scheduler模块中的Request队列中；所以用户可以根据自己的需求进行个性化实现爬虫解析功能。


### [Scheduler](http://godoc.org/github.com/hu17889/go_spider/core/scheduler)

**功能**：用户一般无需自己实现。Scheduler实际上是一个Request对象队列，用来保存尚未被爬取的页面链接和相应的信息，当前队列是缓存到内存中（QueueScheduler），后续会增加基于Redis的队列，解决Spider异常失败后未爬取链接丢失问题；


### [Pipeline](http://godoc.org/github.com/hu17889/go_spider/core/pipeline)

**功能**：用户可以选择自己实现。此模块主要完成数据的输出与持久化。在PageProcesser模块中可用数据被存入了Page对象中的PageItems对象中，此处会获取PageItems的结果并按照自己的要求输出。已有的样例有：PipelineConsole（输出到标准输出），PipelineFile（输出到文件中）
