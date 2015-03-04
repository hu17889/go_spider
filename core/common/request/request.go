// Package request implements request entity contains url and other relevant informaion.
package request

import (
    "net/http"
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

func (this *Request) GetExtension() interface{} {
    return this.extension
}
