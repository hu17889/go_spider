// Package request implements request entity contains url and other relevant informaion.
package request

import (
	"io/ioutil"
	"net/http"
	"os"
	"github.com/bitly/go-simplejson"

	"github.com/hu17889/go_spider/core/common/mlog"
)

// Request represents object waiting for being crawled.
type Request struct {
    url string

    // Responce type: html json jsonp text
    respType string

    // GET POST
    method string

    // POST data
    postdata string

    // name for marking url and distinguish different urls in PageProcesser and Pipeline
    urltag string

    // http header
    header http.Header

    // http cookies
    cookies []*http.Cookie

    // Redirect function for downloader used in http.Client
    // If CheckRedirect returns an error, the Client's Get
    // method returns both the previous Response.
    // If CheckRedirect returns error.New("normal"), the error process after client.Do will ignore the error.
    checkRedirect func(req *http.Request, via []*http.Request) error

    meta interface{}
}

// NewRequest returns initialized Request object.
// The respType is json, jsonp, html, text
/*
func NewRequestSimple(url string, respType string, urltag string) *Request {
    return &Request{url:url, respType:respType}
}
*/

func NewRequest(url string, respType string, urltag string, method string,
    postdata string, header http.Header, cookies []*http.Cookie,
    checkRedirect func(req *http.Request, via []*http.Request) error,
    meta interface{}) *Request {
    return &Request{url, respType, method, postdata, urltag, header, cookies, checkRedirect, meta}
}

func NewRequestWithHeaderFile(url string, respType string, headerFile string) *Request {
	_, err := os.Stat(headerFile)
	if err != nil {
		//file is not exist , using default mode
		return NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	}
	//read file , parse the header and cookies
	b, err := ioutil.ReadFile(headerFile)
	if err != nil {
		//make be:  share access error
		mlog.LogInst().LogError(err.Error())
	}
	js, _ := simplejson.NewJson(b)
	//constructed to header

	h := make(http.Header)
	h.Add("User-Agent", js.Get("User-Agent").MustString())
	h.Add("Cookie", js.Get("Cookie").MustString())
	h.Add("Cache-Control", "max-age=0")
	h.Add("Connection", "keep-alive")
	return NewRequest(url, respType, "", "GET", "", h, nil, nil, nil)
}

func (this *Request) GetUrl() string {
    return this.url
}

func (this *Request) GetUrlTag() string {
    return this.urltag
}

func (this *Request) GetMethod() string {
    return this.method
}

func (this *Request) GetPostdata() string {
    return this.postdata
}

func (this *Request) GetHeader() http.Header {
    return this.header
}

func (this *Request) GetCookies() []*http.Cookie {
    return this.cookies
}

func (this *Request) GetResponceType() string {
    return this.respType
}

func (this *Request) GetRedirectFunc() func(req *http.Request, via []*http.Request) error {
    return this.checkRedirect
}

func (this *Request) GetMeta() interface{} {
    return this.meta
}
