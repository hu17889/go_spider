// craw master module
package spider

import (
    "core/common/mcounter"
    "core/common/mlog"
    "core/common/page"
    "core/common/request"
    "core/downloader"
    "core/page_processer"
    "core/scheduler"
    "fmt"
    //"time"
    //"math/rand"
)

type Spider struct {
    //pStrace *mlog.Strace

    // 页面分析模块
    pPageProcesser page_processer.PageProcesser

    // 下载模块
    pDownloader downloader.Downloader

    // 任务模块
    pScheduler scheduler.Scheduler

    // 并发计数器
    mc  mcounter.Mcounter

    // 并发个数
    threadnum int

    exitWhenComplete bool
}

func NewSpider(pageinst page_processer.PageProcesser) *Spider {
    ap := &Spider{pPageProcesser: pageinst}

    // 初始化
    if ap.pScheduler == nil {
        var s scheduler.Scheduler
        s = scheduler.NewQueueScheduler()
        ap.SetScheduler(s)
    }

    if ap.pDownloader == nil {
        var d downloader.Downloader
        d = downloader.NewHttpDownloader()
        ap.SetDownloader(d)
    }

    ap.exitWhenComplete = true

    return ap
}

// 处理单个url,完成即退出
func (this *Spider) Get(url string, respType string) {
    var urls []string
    urls = append(urls, url)
    this.GetAll(urls, respType)
}

// 批处理url,完成即退出
func (this *Spider) GetAll(urls []string, respType string) {
    // push url
    for _, u := range urls {
        requ := request.NewRequest(u, respType)
        this.addRequest(requ)
    }

    this.Run()
}

func (this *Spider) Run() {
    if this.threadnum <= 0 {
        this.threadnum = 1
    }
    this.mc = mcounter.NewStaticMcounter(this.threadnum)

    for {
        var requ *request.Request
        requ = this.pScheduler.Poll()

        // mc is not atomic
        if this.mc.Count() == 0 && requ == nil && this.exitWhenComplete {
            break
        } else if requ == nil {
            continue
        }
        this.mc.Incr()

        // Asynchronous fetching
        go func(*request.Request) {
            defer func() {
                if r := recover(); r != nil {
                    errStr := fmt.Sprintf("%v", r)
                    mlog.Filelog.LogError("down error " + errStr)
                }
            }()
            defer this.mc.Decr()
            //time.Sleep( time.Duration(rand.Intn(5)) * time.Second)
            this.pageProcess(requ)
        }(requ)
    }
}

// 任务模块
func (this *Spider) SetScheduler(s scheduler.Scheduler) *Spider {
    this.pScheduler = s
    return this
}

func (this *Spider) GetScheduler() scheduler.Scheduler {
    return this.pScheduler
}

// 下载模块
func (this *Spider) SetDownloader(d downloader.Downloader) *Spider {
    this.pDownloader = d
    return this
}

func (this *Spider) GetDownloader() downloader.Downloader {
    return this.pDownloader
}

// 设置处理单元的并发个数
func (this *Spider) SetThreadnum(i int) *Spider {
    this.threadnum = i
    return this
}

func (this *Spider) GetThreadnum() int {
    return this.threadnum
}

// 完成后是否退出程序;
// 当任务队列存在外部插入的时候(如：redis队列)，可以设置为false;
// 当使用QueueScheduler的时候，请设置为默认值true;
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
        mlog.Filelog.LogError("request is nil")
        return
    } else if req.GetUrl() == "" {
        mlog.Filelog.LogError("request is empty")
        return
    }
    this.pScheduler.Push(req)
}

// core processer
func (this *Spider) pageProcess(req *request.Request) {
    //mlog.Filelog.LogError("test")
    var p *page.Page
    p = this.pDownloader.Download(req)
    if p == nil {
        return
    }
    //this.pDownloader.Download(req)
    //fmt.Printf("%v\n", p)

    // todo 下载重试

    this.pPageProcesser.Process(p)
    for _, req := range p.GetTargetRequests() {
        this.addRequest(req)
    }
    // todo pipeline

}
