go_spider
=========
[![Build Status](https://travis-ci.org/hu17889/go_spider.svg)](https://travis-ci.org/hu17889/go_spider)


A crawler of vertical communities that achieved by GOLANG. 

![image](https://github.com/hu17889/doc/blob/master/go_spider/img/logo.png)


Latest stable Release: [Version 1.0 (Sep 23, 2014)](https://github.com/hu17889/go_spider/releases).

* QQ群号：337344607


## Features

* Concurrent 
* Suit for vertical communities
* Flexible, Modularization
* Native Go implementation
* Can be expanded to individualization easily


## Requirements

* Go 1.1 or higher

## Documentation

[中文文档](https://github.com/hu17889/go_spider/wiki/%E4%B8%AD%E6%96%87%E6%96%87%E6%A1%A3) && [常见问题](https://github.com/hu17889/go_spider/wiki/%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98%E4%B8%8E%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E).


## Installation

```
go get github.com/hu17889/go_spider
go get github.com/PuerkitoBio/goquery
go get github.com/bitly/go-simplejson
```

This project is dependent on [simplejson](https://github.com/bitly/go-simplejson/blob/master/simplejson.go), [goquery](https://github.com/PuerkitoBio/goquery).


## Use example

Here is an example for crawl github content. You can have a try for experience the crawl process.
* `go install github.com/hu17889/go_spider/example/github_repo_page_processor`
* `./bin/github_repo_page_processor`

More examples here: [examples](https://github.com/hu17889/go_spider/tree/master/example).


## Make your spider

``` Go
    // Spider input:
    //  PageProcesser ;
    //  Task name used in Pipeline for record;
    spider.NewSpider(NewMyPageProcesser(), "TaskName").
        AddUrl("https://github.com/hu17889?tab=repositories", "html"). // Start url, html is the responce type ("html" or "json")
        AddPipeline(pipeline.NewPipelineConsole()).                    // Print result on screen
        SetThreadnum(3).                                               // Crawl request by three Coroutines
        Run()
```

- Use default modules 

 - Downloader：HttpDownloader
 - Scheduler：QueueScheduler
 - Pipeline：PipelineConsole，PipelineFile

- Use your modules

Just copy the default modules and modify it!

If you make a Downloader module, you can use it by `Spider.SetDownloader(your_downloader)`.

If you make a Pipeline module, you can use it by `Spider.AddPipeline(your_pipeline)`.

If you make a Scheduler module, you can use it by `Spider.SetScheduler(your_scheduler)`.


## License
go_spider is licensed under the [Mozilla Public License Version 2.0](https://github.com/hu17889/go_spider/blob/master/LICENSE)

Mozilla summarizes the license scope as follows:
> MPL: The copyleft applies to any files containing MPLed code.


That means:
  * You can **use** the **unchanged** source code both in private as also commercial
  * You **needn't publish** the source code of your library as long the files licensed under the MPL 2.0 are **unchanged**
  * You **must publish** the source code of any **changed files** licensed under the MPL 2.0 under a) the MPL 2.0 itself or b) a compatible license (e.g. GPL 3.0 or Apache License 2.0)

Please read the [MPL 2.0 FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html) if you have further questions regarding the license.

You can read the full terms here: [LICENSE](https://raw.github.com/go-sql-driver/mysql/master/LICENSE).
