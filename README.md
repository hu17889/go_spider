go_spider
=========
[![Build Status](https://travis-ci.org/hu17889/go_spider.svg)](https://travis-ci.org/hu17889/go_spider)


A crawler of vertical communities that achieved by GOLANG. 

![image](https://github.com/hu17889/doc/blob/master/go_spider/img/logo.png)


Latest stable Release: [Version 1.0 (Sep 23, 2014)](https://github.com/hu17889/go_spider/releases).


## Features


## Requirements


## Documentation

[中文文档](https://github.com/hu17889/go_spider/wiki/%E4%B8%AD%E6%96%87%E6%96%87%E6%A1%A3).


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


## Useage


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
