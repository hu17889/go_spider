// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package page contain result catched by download.
// And it alse has result parsed by PageProcesser.
package page

import (
    "core/common/request"
    //"core/common/page_items"
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
    // pItems *page_items.PageItems

    // html 字符串
    body string

    // 待抓取目标Request
    targetRequests []*request.Request
}

func NewPage() *Page {
    return &Page{}
}

// request struct
func (this *Page) SetRequest(r *request.Request) *Page {
    this.req = r
    return this
}

func (this *Page) GetRequest() *request.Request {
    return this.req
}

// 加入新的待爬取Request
// 爬取的Request的返回类型和当前Request相同
func (this *Page) AddTargetRequest(s string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequest(s, this.req.GetResponceType()))
    return this
}

// 获取目标Request，会在Spider中调用并将目标Request插入Scheduler
func (this *Page) getTargetRequests() []*request.Request {
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
