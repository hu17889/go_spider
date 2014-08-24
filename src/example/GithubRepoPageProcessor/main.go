//
package main

import (
    "core/common/page"
    "core/page_processer"
    "core/spider"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "strings"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// we parse html dom here and get the content that we want.
// we use goquery (http://godoc.org/github.com/PuerkitoBio/goquery#Selection.Html) to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
    query := p.GetHtmlParser()
    var urls []string
    query.Find("a[class='css-truncate css-truncate-target']").Each(func(i int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        urls = append(urls, "http://github.com/"+href)
        fmt.Printf("%v\n", urls)
    })
    // these urls will be saved and crawed by other Coroutines
    p.AddTargetRequests(urls, "html")

    name := query.Find(".entry-title .author").Text()
    name = strings.Trim(name, " \t\n")
    repository := query.Find(".entry-title .js-current-repository").Text()
    repository = strings.Trim(repository, " \t\n")
    if name == "" {
        p.SetSkip(true)
    }
    // the entity we want to save by Pipeline
    p.AddField("name", name+"/"+repository)
}

func main() {
    var pagepro page_processer.PageProcesser
    pagepro = NewMyPageProcesser()
    spider.NewSpider(pagepro).
        AddUrl("https://github.com/hu17889?tab=repositories", "html"). // start url
        SetThreadnum(3).                                               // craw reques in three Coroutines
        Run()
}
