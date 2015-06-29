package main

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/pipeline"
    "github.com/hu17889/go_spider/core/spider"
    "net/http"
    "net/url"
    "strings"
    "fmt"
    "errors"
)

type MyPageProcesser struct {
    cookies []*http.Cookie
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {

    if p.GetUrlTag() == "site_login" {
        //fmt.Printf("%v\n", p.GetCookies())
        this.cookies = p.GetCookies()
        // AddTargetRequestWithParams Params:
        //  1. Url.
        //  2. Responce type is "html" or "json" or "jsonp" or "text".
        //  3. The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
        //  4. The method is POST or GET.
        //  5. The postdata is body string sent to sever.
        //  6. The header is header for http request.
        //  7. Cookies
        //  8. Http redirect function
        if len(this.cookies) != 0 {
            p.AddField("info", "get cookies success")
            req := request.NewRequest("http://backadmin.hucong.net/site/index", "html", "site_index", "GET", "", nil, this.cookies, nil, nil)
            p.AddTargetRequestWithParams(req)
        } else {
            p.AddField("info", "get cookies failed")
        }
    } else {
        //fmt.Printf("%v\n", p.GetBodyStr())
        query := p.GetHtmlParser()
        pageTitle := query.Find(".page-content .page-title").Text()

        if len(pageTitle) != 0 {
            p.AddField("page_title", pageTitle)
            p.AddField("info", "login success")
        } else {
            p.AddField("info", "login failed")
        }

    }

    return
    if !p.IsSucc() {
        println(p.Errormsg())
        return
    }

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

func (this *MyPageProcesser) Finish() {
    fmt.Printf("TODO:before end spider \r\n")
}

// function that prevent redirect for getting cookies
// If CheckRedirect function returns error.New("normal"), the error process after client.Do will ignore the error.
func myRedirect(req *http.Request, via []*http.Request) error {
    return errors.New("normal")
}

func main() {

    // POST data
    post_arg := url.Values{
        "name": {"admin"},
        "pwd":  {"admin"},
    }

    // http header
    header := make(http.Header)
    header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Spider input:
    //  PageProcesser ;
    //  Task name used in Pipeline for record;
    // AddRequest Params:
    //  1. Url.
    //  2. Responce type is "html" or "json" or "jsonp" or "text".
    //  3. The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
    //  4. The method is POST or GET.
    //  5. The postdata is body string sent to sever.
    //  6. The header is header for http request.
    //  7. Cookies
    //  8. Http redirect function
    req := request.NewRequest("http://backadmin.hucong.net/main/user/login", "html", "site_login", "POST", post_arg.Encode(), header, nil, myRedirect, nil)

    spider.NewSpider(NewMyPageProcesser(), "TaskName").
        AddRequest(req).
        AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
        SetThreadnum(3).                            // Crawl request by three Coroutines
        Run()
}
