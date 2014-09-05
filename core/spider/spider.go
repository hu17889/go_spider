// craw master module
package spider

import (
    "fmt"
    "github.com/hu17889/go_spider/core/common/etc_config"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page_items"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/common/resource_manage"
    "github.com/hu17889/go_spider/core/downloader"
    "github.com/hu17889/go_spider/core/page_processer"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/scheduler"
    //"time"
    //"math/rand"
)

type Spider struct {
    debug    bool // make panic effective
    taskname string

    pPageProcesser page_processer.PageProcesser

    pDownloader downloader.Downloader

    pScheduler scheduler.Scheduler

    pPiplelines []pipeline.Pipeline

    mc  resource_manage.ResourceManage

    threadnum int

    exitWhenComplete bool
}

// Spider is scheduler module for all the other modules, like downloader, pipeline, scheduler and etc.
// The confpath could be empty string, then config will use default path "WD/etc/main.conf"
// The taskname could be empty string too, or it can be used in Pipeline for record the result crawled by which task;
func NewSpider(pageinst page_processer.PageProcesser, confpath string, taskname string) *Spider {
    // init config
    etc_config.StartConf(confpath)

    ap := &Spider{taskname: taskname, pPageProcesser: pageinst}

    ap.debug = false
    ap.exitWhenComplete = true

    // init spider
    if ap.pScheduler == nil {
        ap.SetScheduler(scheduler.NewQueueScheduler())
    }

    if ap.pDownloader == nil {
        ap.SetDownloader(downloader.NewHttpDownloader())
    }

    ap.pPiplelines = make([]pipeline.Pipeline, 0)

    return ap
}

func (this *Spider) SetDebug(debug bool) *Spider {
    this.debug = debug
    return this
}

func (this *Spider) Taskname() string {
    return this.taskname
}

// Deal with one url and return the PageItems
func (this *Spider) Get(url string, respType string) *page_items.PageItems {
    var urls []string
    urls = append(urls, url)
    items := this.GetAll(urls, respType)
    if len(items) != 0 {
        return items[0]
    }
    return nil
}

// Deal with several urls and return the PageItems slice
func (this *Spider) GetAll(urls []string, respType string) []*page_items.PageItems {
    // push url
    for _, u := range urls {
        req := request.NewRequest(u, respType)
        this.addRequest(req)
    }

    pip := pipeline.NewCollectPipelinePageItems()
    this.AddPipeline(pip)

    this.Run()

    return pip.GetCollected()
}

func (this *Spider) Run() {
    if this.threadnum <= 0 {
        this.threadnum = 1
    }
    this.mc = resource_manage.NewResourceManageChan(this.threadnum)

    for {
        req := this.pScheduler.Poll()

        // mc is not atomic
        if this.mc.Has() == 0 && req == nil && this.exitWhenComplete {
            break
        } else if req == nil {
            continue
        }
        this.mc.GetOne()

        // Asynchronous fetching
        go func(*request.Request) {
            if !this.debug {
                defer func() {
                    if r := recover(); r != nil {
                        errStr := fmt.Sprintf("%v", r)
                        mlog.LogInst().LogError("down error: " + errStr)
                    }
                }()
            }
            defer this.mc.FreeOne()
            //time.Sleep( time.Duration(rand.Intn(5)) * time.Second)
            this.pageProcess(req)
        }(req)
    }
    this.close()
}

func (this *Spider) close() {
    this.debug = false
    this.SetScheduler(scheduler.NewQueueScheduler())
    this.SetDownloader(downloader.NewHttpDownloader())
    this.pPiplelines = make([]pipeline.Pipeline, 0)
    this.exitWhenComplete = true
}

func (this *Spider) AddPipeline(p pipeline.Pipeline) *Spider {
    this.pPiplelines = append(this.pPiplelines, p)
    return this
}

func (this *Spider) SetScheduler(s scheduler.Scheduler) *Spider {
    this.pScheduler = s
    return this
}

func (this *Spider) GetScheduler() scheduler.Scheduler {
    return this.pScheduler
}

func (this *Spider) SetDownloader(d downloader.Downloader) *Spider {
    this.pDownloader = d
    return this
}

func (this *Spider) GetDownloader() downloader.Downloader {
    return this.pDownloader
}

func (this *Spider) SetThreadnum(i int) *Spider {
    this.threadnum = i
    return this
}

func (this *Spider) GetThreadnum() int {
    return this.threadnum
}

// If exit when each crawl task is done.
// If you want to keep spider in memory all the time and add url from outside, you can set it true.
func (this *Spider) SetExitWhenComplete(e bool) *Spider {
    this.exitWhenComplete = e
    return this
}

func (this *Spider) GetExitWhenComplete() bool {
    return this.exitWhenComplete
}

func (this *Spider) AddUrl(url string, respType string) *Spider {
    req := request.NewRequest(url, respType)
    this.addRequest(req)
    return this
}

func (this *Spider) AddUrls(urls []string, respType string) *Spider {
    for _, url := range urls {
        req := request.NewRequest(url, respType)
        this.addRequest(req)
    }
    return this
}

// add Request to Schedule
func (this *Spider) addRequest(req *request.Request) {
    if req == nil {
        mlog.LogInst().LogError("request is nil")
        return
    } else if req.GetUrl() == "" {
        mlog.LogInst().LogError("request is empty")
        return
    }
    this.pScheduler.Push(req)
}

// core processer
func (this *Spider) pageProcess(req *request.Request) {
    p := this.pDownloader.Download(req)
    if p == nil {
        return
    }

    // TODO: download retry

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

}
