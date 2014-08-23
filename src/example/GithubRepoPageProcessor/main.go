// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package main

import (
    "core/common/page"
    "core/page_processer"
    "core/spider"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

// 解析page内容，提取a标签等,用来
func (this *MyPageProcesser) Process(p *page.Page) {
}

func main() {
    var pagepro page_processer.PageProcesser
    pagepro = NewMyPageProcesser()
    //var sp *spider.Spider
    spider.NewSpider(pagepro).Get("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1", "json")

    //sp.Get("http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1", "json")

    /*
    	urls := []string{
    		"http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1",
    		"http://live.sina.com.cn/zt/api/l/get/finance/globalnews1/index.htm?format=json&id=23521&pagesize=4&dire=f&dpc=1",
    	}
    	sp.GetAll(urls, "json")
    */
}
