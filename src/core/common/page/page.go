// Package page contains result catched by Downloader.
// And it alse has result parsed by PageProcesser.
package page

import (
    "core/common/page_items"
    "core/common/request"
    "github.com/PuerkitoBio/goquery"
)


// Page represents an entity be crawled.
type Page struct {
    // The request is crawled by spider that contains url and relevent information.
    req *request.Request

    // The body is plain text of crawl result.
    body string

    // The docParser is a pointer of goquery boject that contains html result.
    docParser *goquery.Document

    // The jsonMap is the json result.
    jsonMap interface{}

    // The pItems is object for save Key-Values in PageProcesser.
    // And pItems is output in Pipline.
    pItems *page_items.PageItems

    // The targetRequests is requests to put into Scheduler.
    targetRequests []*request.Request
}

// NewPage returns initialized Page object.
func NewPage(req *request.Request) *Page {
    return &Page{pItems: page_items.NewPageItems(req), req: req}
}

// AddField saves KV pair to PageItems preparing for Pipeline
func (this *Page) AddField(name string, value string) {
    this.pItems.AddItem(name, value)
}

// GetPageItems returns PageItems object that record KV pair parsed in PageProcesser.
func (this *Page) GetPageItems() *page_items.PageItems {
    return this.pItems
}

// SetSkip set label "skip" of PageItems.
// PageItems will not be saved in Pipeline wher skip is set true
func (this *Page) SetSkip(skip bool) {
    this.pItems.SetSkip(skip)
}

// GetSkip returns skip label of PageItems.
func (this *Page) GetSkip() bool {
    return this.pItems.GetSkip()
}

// SetRequest saves request oject of this page.
func (this *Page) SetRequest(r *request.Request) *Page {
    this.req = r
    return this
}

// GetRequest returns request oject of this page.
func (this *Page) GetRequest() *request.Request {
    return this.req
}

// AddTargetRequest adds new Request waitting for crawl.
func (this *Page) AddTargetRequest(url string, respType string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType))
    return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
    for _, url := range urls {
        this.AddTargetRequest(url, respType)
    }
    return this
}

// GetTargetRequests returns the target requests that will put into Scheduler
func (this *Page) GetTargetRequests() []*request.Request {
    return this.targetRequests
}

// SetBodyStr saves plain string crawled in Page.
func (this *Page) SetBodyStr(body string) *Page {
    this.body = body
    return this
}

// GetBodyStr returns plain string crawled.
func (this *Page) GetBodyStr() string {
    return this.body
}

// SetHtmlParser saves goquery object binded to target crawl result.
func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
    this.docParser = doc
    return this
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (this *Page) GetHtmlParser() *goquery.Document {
    return this.docParser
}

// SetJsonMap saves json result. 
func (this *Page) SetJsonMap(str interface{}) *Page {
    this.jsonMap = str
    return this
}

// SetJsonMap returns json result. 
func (this *Page) GetJsonMap() interface{} {
    return this.jsonMap
}
