// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package downloader

import (
    "core/common/mlog"
    "core/common/page"
    "core/common/request"
    "encoding/json"
    "github.com/PuerkitoBio/goquery"
    "io/ioutil"
    "net/http"
)

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
        mlog.Log("error request type:" + mtype)
        return nil
    }
}

func (this *HttpDownloader) downloadHtml(req *request.Request) *page.Page {
    var err error
    var url string
    url = req.GetUrl()

    var resp *http.Response
    resp, err = http.Get(url)
    if err != nil {
        mlog.Log(err)
    }
    defer resp.Body.Close()

    var doc *goquery.Document
    if doc, err = goquery.NewDocumentFromReader(resp.Body); err != nil {
        mlog.Log(err)
    }

    var body string
    body, _ = doc.Html()

    // create Page
    var p *page.Page = page.NewPage().
        SetRequest(req).
        SetBodyStr(body).
        SetHtmlParser(doc)

    return p

}

func (this *HttpDownloader) downloadJson(req *request.Request) *page.Page {
    var err error
    var url string
    url = req.GetUrl()

    var resp *http.Response
    resp, err = http.Get(url)
    if err != nil {
        mlog.Log(err)
    }
    defer resp.Body.Close()

    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        mlog.Log(err)
    }

    var r interface{}
    err = json.Unmarshal(body, &r)
    if err != nil {
        mlog.Log(err)
    }

    // create Page
    // json result
    var p *page.Page = page.NewPage().
        SetRequest(req).
        SetBodyStr(string(body)).
        SetJsonMap(r)

    return p
}
