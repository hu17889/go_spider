package downloader

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/common/request"
    "io/ioutil"
    "net/http"
    //"fmt"
    "strings"
)

// The HttpDownloader download page by http.
// Html content is contained in dom parser of package goquery.
// Json content is saved
// The page result is saved in Page and
type HttpDownloader struct {
}

func NewHttpDownloader() *HttpDownloader {
    return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *request.Request) *page.Page {
    var mtype string
    mtype = req.GetResponceType()
    switch mtype {
    case "html":
        return this.downloadHtml(req)
    case "json":
        return this.downloadJson(req)
    default:
        mlog.LogInst().LogError("error request type:" + mtype)
        return nil
    }
}

// The acceptableCharset is test for whether Content-Type is UTF-8 or not
func (this *HttpDownloader) acceptableCharset(contentTypes []string) bool {
    // each type is like [text/html; charset=UTF-8]
    // we want the UTF-8 only
    for _, cType := range contentTypes {
        if strings.Index(cType, "UTF-8") != -1 || strings.Index(cType, "utf-8") != -1 {
            return true
        }
    }
    return false
}

func (this *HttpDownloader) downloadHtml(req *request.Request) *page.Page {
    var p *page.Page = page.NewPage(req)
    var err error
    var url string
    if url = req.GetUrl(); len(url) == 0 {
        mlog.LogInst().LogError("url is empty")
        p.SetStatus(true, "url is empty")
        return p
    }

    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }
    defer resp.Body.Close()

    /*
       if ok := this.acceptableCharset(resp.Header["Content-Type"]); !ok {
           mlog.LogInst().LogError(fmt.Sprintf("Content-Type is not UTF-8 : %v",resp.Header["Content-Type"]))
           p.SetStatus(true, fmt.Sprintf("Content-Type is not UTF-8 : %v", resp.Header["Content-Type"]))
           return p
       }
    */

    var doc *goquery.Document
    if doc, err = goquery.NewDocumentFromReader(resp.Body); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    var body string
    if body, err = doc.Html(); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    p.SetBodyStr(body).SetHtmlParser(doc).SetStatus(false, "")

    return p
}

func (this *HttpDownloader) downloadJson(req *request.Request) *page.Page {
    var p *page.Page = page.NewPage(req)
    var err error
    var url string
    if url = req.GetUrl(); len(url) == 0 {
        mlog.LogInst().LogError("url is empty")
        p.SetStatus(true, "url is empty")
        return p
    }

    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }
    defer resp.Body.Close()

    var body []byte
    if body, err = ioutil.ReadAll(resp.Body); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    var r *simplejson.Json
    if r, err = simplejson.NewJson(body); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    // json result
    p.SetBodyStr(string(body)).SetJson(r).SetStatus(false, "")

    return p
}
