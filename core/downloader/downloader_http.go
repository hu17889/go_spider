package downloader

import (
    "bytes"
    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    iconv "github.com/djimenez/iconv-go"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/common/util"
    "io"
    "io/ioutil"
    "net/http"
    "regexp"
    //"golang.org/x/net/html"
    //"fmt"
    "strings"
)

// The HttpDownloader download page by package net/http.
// The "html" content is contained in dom parser of package goquery.
// The "json" content is saved.
// The "jsonp" content is modified to json.
// The "text" content will save body plain text only.
// The page result is saved in Page.
type HttpDownloader struct {
}

func NewHttpDownloader() *HttpDownloader {
    return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *request.Request) *page.Page {
    var mtype string
    var p = page.NewPage(req)
    mtype = req.GetResponceType()
    switch mtype {
    case "html":
        return this.downloadHtml(p, req)
    case "json":
        fallthrough
    case "jsonp":
        return this.downloadJson(p, req)
    case "text":
        return this.downloadText(p, req)
    default:
        mlog.LogInst().LogError("error request type:" + mtype)
    }
    return p
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

// The getCharset used for parsing the header["Content-Type"] string to get charset of the page.
func (this *HttpDownloader) getCharset(header http.Header) string {
    reg, err := regexp.Compile("charset=(.*)$")
    if err != nil {
        mlog.LogInst().LogError(err.Error())
        return ""
    }

    var charset string
    for _, cType := range header["Content-Type"] {
        substrings := reg.FindStringSubmatch(cType)
        if len(substrings) == 2 {
            charset = substrings[1]
        }
    }

    return charset
}

// Use golang.org/x/text/encoding. Get page body and change it to utf-8
func (this *HttpDownloader) changeCharsetEncoding(charset string, sor io.ReadCloser) string {
    ischange := true
    var tr transform.Transformer
    cs := strings.ToLower(charset)
    if cs=="gbk" {
        tr = simplifiedchinese.GBK.NewDecoder()
    } else if cs=="gb18030" {
        tr = simplifiedchinese.GB18030.NewDecoder()
    } else if cs=="hzgb2312" || cs=="gb2312" || cs=="hz-gb2312" {
        tr = simplifiedchinese.HZGB2312.NewDecoder()
    } else {
        ischange = false
    }

    var destReader io.Reader
    if ischange {
        transReader := transform.NewReader(sor, tr)
        destReader = transReader
    } else {
        destReader = sor
    }
        
    var sorbody []byte
    var err error
    if sorbody, err = ioutil.ReadAll(destReader); err != nil {
        mlog.LogInst().LogError(err.Error())
        return ""
    }
    bodystr := string(sorbody)


    return bodystr
}

// Use go-iconv. Get page body and change it to utf-8
func (this *HttpDownloader) changeCharsetGoIconv(charset string, sor io.ReadCloser) string {
    var err error
    var converter *iconv.Converter
    if charset != "" && strings.ToLower(charset) != "utf-8" && strings.ToLower(charset) != "utf8" {
        converter, err = iconv.NewConverter(charset, "utf-8")
        if err != nil {
            mlog.LogInst().LogError(err.Error())
            return ""
        }
        defer converter.Close()
    }

    var sorbody []byte
    if sorbody, err = ioutil.ReadAll(sor); err != nil {
        mlog.LogInst().LogError(err.Error())
        return ""
    }
    bodystr := string(sorbody)

    var destbody string
    if converter != nil {
        // convert to utf8
        destbody, err = converter.ConvertString(bodystr)
        if err != nil {
            mlog.LogInst().LogError(err.Error())
            return ""
        }
    } else {
        destbody = bodystr
    }
    return destbody
}


// Download file and change the charset of page charset.
func (this *HttpDownloader) downloadFile(p *page.Page, req *request.Request) (*page.Page, string) {
    var err error
    var url string
    if url = req.GetUrl(); len(url) == 0 {
        mlog.LogInst().LogError("url is empty")
        p.SetStatus(true, "url is empty")
        return p, ""
    }

    var resp *http.Response
    if resp, err = http.Get(url); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p, ""
    }
    p.SetHeader(resp.Header)
    p.SetCookies(resp.Cookies())

    // get converter to utf-8
    charset := this.getCharset(resp.Header)

    bodyStr := this.changeCharsetEncoding(charset, resp.Body)
    defer resp.Body.Close()
    return p, bodyStr
}

func (this *HttpDownloader) downloadHtml(p *page.Page, req *request.Request) *page.Page {
    var err error
    p, destbody := this.downloadFile(p, req)
    if !p.IsSucc() {
        return p
    }
    bodyReader := bytes.NewReader([]byte(destbody))

    var doc *goquery.Document
    if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
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

func (this *HttpDownloader) downloadJson(p *page.Page, req *request.Request) *page.Page {
    var err error
    p, destbody := this.downloadFile(p, req)
    if !p.IsSucc() {
        return p
    }

    var body []byte
    body = []byte(destbody)
    mtype := req.GetResponceType()
    if mtype == "jsonp" {
        tmpstr := util.JsonpToJson(destbody)
        body = []byte(tmpstr)
    }

    var r *simplejson.Json
    if r, err = simplejson.NewJson(body); err != nil {
        mlog.LogInst().LogError(string(body) + "\t" + err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    // json result
    p.SetBodyStr(string(body)).SetJson(r).SetStatus(false, "")

    return p
}

func (this *HttpDownloader) downloadText(p *page.Page, req *request.Request) *page.Page {
    p, destbody := this.downloadFile(p, req)
    if !p.IsSucc() {
        return p
    }

    p.SetBodyStr(destbody).SetStatus(false, "")
    return p
}
