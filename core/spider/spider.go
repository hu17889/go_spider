// craw master module
package spider

import (
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/common/page_items"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/common/resource_manage"
    "github.com/hu17889/go_spider/core/downloader"
    "github.com/hu17889/go_spider/core/page_processer"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/scheduler"
    "math/rand"
    //"net/http"
    "time"
    //"fmt"
)

type Spider struct {
    taskname string

    pPageProcesser page_processer.PageProcesser

    pDownloader downloader.Downloader

    pScheduler scheduler.Scheduler

    pPiplelines []pipeline.Pipeline

    mc  resource_manage.ResourceManage

    threadnum uint

    exitWhenComplete bool

    // Sleeptype can be fixed or rand.
    startSleeptime uint
    endSleeptime   uint
    sleeptype      string
}

// Spider is scheduler module for all the other modules, like downloader, pipeline, scheduler and etc.
// The taskname could be empty string too, or it can be used in Pipeline for record the result crawled by which task;
func NewSpider(pageinst page_processer.PageProcesser, taskname string) *Spider {
    mlog.StraceInst().Open()

    ap := &Spider{taskname: taskname, pPageProcesser: pageinst}

    // init filelog.
    ap.CloseFileLog()
    ap.exitWhenComplete = true
    ap.sleeptype = "fixed"
    ap.startSleeptime = 0

    // init spider
    if ap.pScheduler == nil {
        ap.SetScheduler(scheduler.NewQueueScheduler(false))
    }

    if ap.pDownloader == nil {
        ap.SetDownloader(downloader.NewHttpDownloader())
    }

    mlog.StraceInst().Println("** start spider **")
    ap.pPiplelines = make([]pipeline.Pipeline, 0)

    return ap
}

func (this *Spider) Taskname() string {
    return this.taskname
}

// Deal with one url and return the PageItems.
func (this *Spider) Get(url string, respType string) *page_items.PageItems {
    req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
    return this.GetByRequest(req)
}

// Deal with several urls and return the PageItems slice.
func (this *Spider) GetAll(urls []string, respType string) []*page_items.PageItems {
    for _, u := range urls {
        req := request.NewRequest(u, respType, "", "GET", "", nil, nil, nil, nil)
        this.AddRequest(req)
    }

    pip := pipeline.NewCollectPipelinePageItems()
    this.AddPipeline(pip)

    this.Run()

    return pip.GetCollected()
}

// Deal with one url and return the PageItems with other setting.
func (this *Spider) GetByRequest(req *request.Request) *page_items.PageItems {
    var reqs []*request.Request
    reqs = append(reqs, req)
    items := this.GetAllByRequest(reqs)
    if len(items) != 0 {
        return items[0]
    }
    return nil
}

// Deal with several urls and return the PageItems slice
func (this *Spider) GetAllByRequest(reqs []*request.Request) []*page_items.PageItems {
    // push url
    for _, req := range reqs {
        //req := request.NewRequest(u, respType, urltag, method, postdata, header, cookies)
        this.AddRequest(req)
    }

    pip := pipeline.NewCollectPipelinePageItems()
    this.AddPipeline(pip)

    this.Run()

    return pip.GetCollected()
}

func (this *Spider) Run() {
    if this.threadnum == 0 {
        this.threadnum = 1
    }
    this.mc = resource_manage.NewResourceManageChan(this.threadnum)

    //init db  by sorawa

    for {
        req := this.pScheduler.Poll()

        // mc is not atomic
        if this.mc.Has() == 0 && req == nil && this.exitWhenComplete {
	    mlog.StraceInst().Println("** executed callback **")
	    this.pPageProcesser.Finish()
            mlog.StraceInst().Println("** end spider **")
            break
        } else if req == nil {
            time.Sleep(500 * time.Millisecond)
            //mlog.StraceInst().Println("scheduler is empty")
            continue
        }
        this.mc.GetOne()

        // Asynchronous fetching
        go func(req *request.Request) {
            defer this.mc.FreeOne()
            //time.Sleep( time.Duration(rand.Intn(5)) * time.Second)
            mlog.StraceInst().Println("start crawl : " + req.GetUrl())
            this.pageProcess(req)
        }(req)
    }
    this.close()
}

func (this *Spider) close() {
    this.SetScheduler(scheduler.NewQueueScheduler(false))
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

func (this *Spider) SetThreadnum(i uint) *Spider {
    this.threadnum = i
    return this
}

func (this *Spider) GetThreadnum() uint {
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

// The OpenFileLog initialize the log path and open log.
// If log is opened, error info or other useful info in spider will be logged in file of the filepath.
// Log command is mlog.LogInst().LogError("info") or mlog.LogInst().LogInfo("info").
// Spider's default log is closed.
// The filepath is absolute path.
func (this *Spider) OpenFileLog(filePath string) *Spider {
    mlog.InitFilelog(true, filePath)
    return this
}

// OpenFileLogDefault open file log with default file path like "WD/log/log.2014-9-1".
func (this *Spider) OpenFileLogDefault() *Spider {
    mlog.InitFilelog(true, "")
    return this
}

// The CloseFileLog close file log.
func (this *Spider) CloseFileLog() *Spider {
    mlog.InitFilelog(false, "")
    return this
}

// The OpenStrace open strace that output progress info on the screen.
// Spider's default strace is opened.
func (this *Spider) OpenStrace() *Spider {
    mlog.StraceInst().Open()
    return this
}

// The CloseStrace close strace.
func (this *Spider) CloseStrace() *Spider {
    mlog.StraceInst().Close()
    return this
}

// The SetSleepTime set sleep time after each crawl task.
// The unit is millisecond.
// If sleeptype is "fixed", the s is the sleep time and e is useless.
// If sleeptype is "rand", the sleep time is rand between s and e.
func (this *Spider) SetSleepTime(sleeptype string, s uint, e uint) *Spider {
    this.sleeptype = sleeptype
    this.startSleeptime = s
    this.endSleeptime = e
    if this.sleeptype == "rand" && this.startSleeptime >= this.endSleeptime {
        panic("startSleeptime must smaller than endSleeptime")
    }
    return this
}

func (this *Spider) sleep() {
    if this.sleeptype == "fixed" {
        time.Sleep(time.Duration(this.startSleeptime) * time.Millisecond)
    } else if this.sleeptype == "rand" {
        sleeptime := rand.Intn(int(this.endSleeptime-this.startSleeptime)) + int(this.startSleeptime)
        time.Sleep(time.Duration(sleeptime) * time.Millisecond)
    }
}

func (this *Spider) AddUrl(url string, respType string) *Spider {
    req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
    this.AddRequest(req)
    return this
}

func (this *Spider) AddUrlEx(url string, respType string, headerFile string, proxyHost string) *Spider {
    req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
    this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
    return this
}

func (this *Spider) AddUrlWithHeaderFile(url string, respType string, headerFile string) *Spider {
    req := request.NewRequestWithHeaderFile(url, respType, headerFile)
    this.AddRequest(req)
    return this
}

func (this *Spider) AddUrls(urls []string, respType string) *Spider {
    for _, url := range urls {
        req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
        this.AddRequest(req)
    }
    return this
}

func (this *Spider) AddUrlsWithHeaderFile(urls []string, respType string, headerFile string) *Spider {
    for _, url := range urls {
        req := request.NewRequestWithHeaderFile(url, respType, headerFile)
        this.AddRequest(req)
    }
    return this
}

func (this *Spider) AddUrlsEx(urls []string, respType string, headerFile string, proxyHost string) *Spider {
    for _, url := range urls {
        req := request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
        this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
    }
    return this
}

// add Request to Schedule
func (this *Spider) AddRequest(req *request.Request) *Spider {
    if req == nil {
        mlog.LogInst().LogError("request is nil")
        return this
    } else if req.GetUrl() == "" {
        mlog.LogInst().LogError("request is empty")
        return this
    }
    this.pScheduler.Push(req)
    return this
}

//
func (this *Spider) AddRequests(reqs []*request.Request) *Spider {
    for _, req := range reqs {
        this.AddRequest(req)
    }
    return this
}

// core processer
func (this *Spider) pageProcess(req *request.Request) {
    var p *page.Page

    defer func() {
        if err := recover(); err != nil { // do not affect other
            if strerr, ok := err.(string); ok {
                mlog.LogInst().LogError(strerr)
            } else {
                mlog.LogInst().LogError("pageProcess error")
            }
        }
    }()

    // download page
    for i := 0; i < 3; i++ {
        this.sleep()
        p = this.pDownloader.Download(req)
        if p.IsSucc() { // if fail retry 3 times
            break
        }

    }

    if !p.IsSucc() { // if fail do not need process
        return
    }

    this.pPageProcesser.Process(p)
    for _, req := range p.GetTargetRequests() {
        this.AddRequest(req)
    }

    // output
    if !p.GetSkip() {
        for _, pip := range this.pPiplelines {
            //fmt.Println("%v",p.GetPageItems().GetAll())
            pip.Process(p.GetPageItems(), this)
        }
    }
}
