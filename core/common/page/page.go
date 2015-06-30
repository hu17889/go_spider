// Package page contains result catched by Downloader.
// And it alse has result parsed by PageProcesser.
package page

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page_items"
    "github.com/hu17889/go_spider/core/common/request"
    "net/http"
    "strings"
    //"fmt"
)

// Page represents an entity be crawled.
type Page struct {
    // The isfail is true when crawl process is failed and errormsg is the fail resean.
    isfail   bool
    errormsg string

    // The request is crawled by spider that contains url and relevent information.
    req *request.Request

    // The body is plain text of crawl result.
    body string

    header  http.Header
    cookies []*http.Cookie

    // The docParser is a pointer of goquery boject that contains html result.
    docParser *goquery.Document

    // The jsonMap is the json result.
    jsonMap *simplejson.Json

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

// SetHeader save the header of http responce
func (this *Page) SetHeader(header http.Header) {
    this.header = header
}

// GetHeader returns the header of http responce
func (this *Page) GetHeader() http.Header {
    return this.header
}

// SetHeader save the cookies of http responce
func (this *Page) SetCookies(cookies []*http.Cookie) {
    this.cookies = cookies
}

// GetHeader returns the cookies of http responce
func (this *Page) GetCookies() []*http.Cookie {
    return this.cookies
}

// IsSucc test whether download process success or not.
func (this *Page) IsSucc() bool {
    return !this.isfail
}

// Errormsg show the download error message.
func (this *Page) Errormsg() string {
    return this.errormsg
}

// SetStatus save status info about download process.
func (this *Page) SetStatus(isfail bool, errormsg string) {
    this.isfail = isfail
    this.errormsg = errormsg
}

// AddField saves KV string pair to PageItems preparing for Pipeline
func (this *Page) AddField(key string, value string) {
    this.pItems.AddItem(key, value)
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

// GetUrlTag returns name of url.
func (this *Page) GetUrlTag() string {
    return this.req.GetUrlTag()
}

// AddTargetRequest adds one new Request waitting for crawl.
func (this *Page) AddTargetRequest(url string, respType string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil))
    return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
    for _, url := range urls {
        this.AddTargetRequest(url, respType)
    }
    return this
}

// AddTargetRequestWithProxy adds one new Request waitting for crawl.
func (this *Page) AddTargetRequestWithProxy(url string, respType string, proxyHost string) *Page {

    this.targetRequests = append(this.targetRequests, request.NewRequestWithProxy(url, respType, "", "GET", "", nil, nil, proxyHost, nil, nil))
    return this
}

// AddTargetRequestsWithProxy adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithProxy(urls []string, respType string, proxyHost string) *Page {
    for _, url := range urls {
        this.AddTargetRequestWithProxy(url, respType, proxyHost)
    }
    return this
}

// AddTargetRequest adds one new Request with header file for waitting for crawl.
func (this *Page) AddTargetRequestWithHeaderFile(url string, respType string, headerFile string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequestWithHeaderFile(url, respType, headerFile))
    return this
}

// AddTargetRequest adds one new Request waitting for crawl.
// The respType is "html" or "json" or "jsonp" or "text".
// The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
// The method is POST or GET.
// The postdata is http body string.
// The header is http header.
// The cookies is http cookies.
func (this *Page) AddTargetRequestWithParams(req *request.Request) *Page {
    this.targetRequests = append(this.targetRequests, req)
    return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithParams(reqs []*request.Request) *Page {
    for _, req := range reqs {
        this.AddTargetRequestWithParams(req)
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

// GetHtmlParser returns goquery object binded to target crawl result.
func (this *Page) ResetHtmlParser() *goquery.Document {
    r := strings.NewReader(this.body)
    var err error
    this.docParser, err = goquery.NewDocumentFromReader(r)
    if err != nil {
        mlog.LogInst().LogError(err.Error())
        panic(err.Error())
    }
    return this.docParser
}

// SetJson saves json result.
func (this *Page) SetJson(js *simplejson.Json) *Page {
    this.jsonMap = js
    return this
}

// SetJson returns json result.
func (this *Page) GetJson() *simplejson.Json {
    return this.jsonMap
}
