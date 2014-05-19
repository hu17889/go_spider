package spider

import (
    "downloader"
    "fmt"
    "scheduler"
    "page_processer"
    "common/request"
    "common/page"
    "common/mcounter"
)

type Spider struct {
    // 页面分析模块
    pPageProcesser page_processer.PageProcesser

    // 下载模块
    pDownloader downloader.Downloader

    // 任务模块
    pScheduler scheduler.Scheduler

    // 并发计数器
    mc mcounter.Mcounter

    // 并发个数
    threadnum int

    exitWhenComplete bool
}
    

func NewSpider(pageinst page_processer.PageProcesser) *Spider {
    ap := &Spider{pPageProcesser:pageinst}

    // 初始化
    if ap.pScheduler == nil {
        var s scheduler.Scheduler
        s = scheduler.NewQueueScheduler()
        ap.SetScheduler(s)
    }

    if ap.threadnum<=0 {
        ap.threadnum = 1
    }

    if ap.pDownloader == nil {
        var d downloader.Downloader
        d = downloader.NewHttpDownloader()
        ap.SetDownloader(d)
    }

    ap.mc = mcounter.NewStaticMcounter(ap.threadnum)

    ap.exitWhenComplete = true

    return ap
}

// 处理单个url,完成即退出
func (this *Spider) Get(url string, respType string) {
    var urls []string
    urls = append(urls,url)
    this.GetAll(urls, respType)
}

// 批处理url,完成即退出
func (this *Spider) GetAll(urls []string, respType string) {
    // push url
    for _, u := range(urls) {
        requ := request.NewRequest(u, respType)
        this.addRequest(requ)
    }

    this.Run()
}

func (this *Spider) Run() {
    for {
        var requ *request.Request
        requ = this.pScheduler.Poll()

        if this.mc.Count() == 0 && requ == nil && this.exitWhenComplete {
            break
        } else if requ == nil {
            continue
        }
        this.mc.Incr()

        go func(*request.Request) {
            this.pageProcess(requ)            
            this.mc.Decr()
        }(requ)
    }
}

// 任务模块
func (this *Spider) SetScheduler(s scheduler.Scheduler) {
    this.pScheduler = s
}

func (this *Spider) GetScheduler() (scheduler.Scheduler) {
    return this.pScheduler
}

// 下载模块
func (this *Spider) SetDownloader(d downloader.Downloader) {
    this.pDownloader = d
}

func (this *Spider) GetDownloader() (downloader.Downloader) {
    return this.pDownloader
}

// 设置处理单元的并发个数
func (this *Spider) SetThreadnum(i int) {
    this.threadnum = i
}

func (this *Spider) GetThreadnum() (int) {
    return this.threadnum
}

// 完成后是否退出程序;
// 当任务队列存在外部插入的时候(如：redis队列)，可以设置为false;
// 当使用QueueScheduler的时候，请设置为默认值true;
func (this *Spider) SetExitWhenComplete(e bool) {
    this.exitWhenComplete = e
}

func (this *Spider) GetExitWhenComplete() bool {
    return this.exitWhenComplete
}



func (this *Spider) addRequest(requ *request.Request) {
    this.pScheduler.Push(requ)
}

// 核心处理函数
func (this *Spider) pageProcess(requ *request.Request) {
    var p *page.Page
    p = this.pDownloader.Download(requ)
    fmt.Printf("%v\n",p)

    // todo 下载重试
}
