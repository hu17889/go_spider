// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package page contain result catched by download.
// And it alse has result parsed by PageProcesser.
package page

import (
    "common/request"
    "github.com/PuerkitoBio/goquery"
)

type Page struct {
    // request has url
    req *request.Request

    // for parse html DOM
    docParser *goquery.Document

    jsonMap interface{}

    // source html doc string
    body string
}

func NewPage() *Page {
    return &Page{}
}

// 请求结构体
func (this *Page) SetRequest(r *request.Request) {
    this.req = r
}

func (this *Page) GetRequest() *request.Request {
    return this.req
}

// 抓取结果的字符串结果
func (this *Page) SetBodyStr(body string) {
    this.body = body
}

func (this *Page) GetBodyStr() string {
    return this.body
}

// html解析结果
func (this *Page) SetHtmlParser(doc *goquery.Document) {
    this.docParser = doc
}

func (this *Page) GetHtmlParser() *goquery.Document {
    return this.docParser
}

// json解析结果
func (this *Page) SetJsonMap(str interface{}) {
    this.jsonMap = str
}

func (this *Page) GetJsonMap() interface{} {
    return this.jsonMap
}
