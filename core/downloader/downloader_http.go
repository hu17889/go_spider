package downloader

import (
    "core/common/mlog"
    "core/common/page"
    "core/common/request"
    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "io/ioutil"
    "net/http"
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

func (this *HttpDownloader) downloadHtml(req *request.Request) *page.Page {
    var err error
    var url string
    if url = req.GetUrl(); len(url) == 0 {
        return nil
    }

    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }
    defer resp.Body.Close()

    var doc *goquery.Document
    if doc, err = goquery.NewDocumentFromReader(resp.Body); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }

    var body string
    if body, err = doc.Html(); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }

    // create Page
    var p *page.Page = page.NewPage(req).
        SetBodyStr(body).
        SetHtmlParser(doc)

    return p

}

func (this *HttpDownloader) downloadJson(req *request.Request) *page.Page {
    var err error
    var url string
    if url = req.GetUrl(); len(url) == 0 {
        mlog.LogInst().LogError(err.Error())
        return nil
    }

    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }
    defer resp.Body.Close()

    var body []byte
    if body, err = ioutil.ReadAll(resp.Body); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }

    var r *simplejson.Json
    if r, err = simplejson.NewJson(body); err != nil {
        mlog.LogInst().LogError(err.Error())
        return nil
    }

    // create Page
    // json result
    var p *page.Page = page.NewPage(req).
        SetBodyStr(string(body)).
        SetJsonMap(r)

    return p
}
