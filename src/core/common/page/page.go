// Package page contain result catched by download.
// And it alse has result parsed by PageProcesser.
package page

import (
    "core/common/page_items"
    "core/common/request"
    "github.com/PuerkitoBio/goquery"
)

type Page struct {
    // request has url
    req *request.Request

    // html DOM解析器
    docParser *goquery.Document

    // json 格式结果
    jsonMap interface{}

    // items to save in pipeline
    pItems *page_items.PageItems

    // html 字符串
    body string

    // 待抓取目标Request
    targetRequests []*request.Request
}

func NewPage(req *request.Request) *Page {
    return &Page{pItems: page_items.NewPageItems(req), req: req}
}

// save KV pair to PageItems preparing for Pipeline
func (this *Page) AddField(name string, value string) {
    this.pItems.AddItem(name, value)
}

func (this *Page) GetPageItems() *page_items.PageItems {
    return this.pItems
}

// PageItems will not be saved in Pipeline wher skip is set true
func (this *Page) SetSkip(skip bool) {
    this.pItems.SetSkip(skip)
}

func (this *Page) GetSkip() bool {
    return this.pItems.GetSkip()
}

// request struct
func (this *Page) SetRequest(r *request.Request) *Page {
    this.req = r
    return this
}

func (this *Page) GetRequest() *request.Request {
    return this.req
}

// Add new Request waitting for craw
func (this *Page) AddTargetRequest(url string, respType string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType))
    return this
}

// Add new Requests waitting for craw
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
    for _, url := range urls {
        this.AddTargetRequest(url, respType)
    }
    return this
}

// 获取目标Request，会在Spider中调用并将目标Request插入Scheduler
func (this *Page) GetTargetRequests() []*request.Request {
    return this.targetRequests
}

// 原始 html string
func (this *Page) SetBodyStr(body string) *Page {
    this.body = body
    return this
}

func (this *Page) GetBodyStr() string {
    return this.body
}

// html 解析器实例，Downloader获取了Request之后，用goquery封装获取Html的Dom结果
func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
    this.docParser = doc
    return this
}

func (this *Page) GetHtmlParser() *goquery.Document {
    return this.docParser
}

// json parse result
func (this *Page) SetJsonMap(str interface{}) *Page {
    this.jsonMap = str
    return this
}

func (this *Page) GetJsonMap() interface{} {
    return this.jsonMap
}
