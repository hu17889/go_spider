go_spider
=========
[![Build Status](https://travis-ci.org/hu17889/go_spider.svg)](https://travis-ci.org/hu17889/go_spider)


**v0.1** First Milestone，完成基本框架

## 简介


本项目基于golang开发，是一个垂直领域的爬虫引擎，主要希望能将各个功能模块区分开，方便使用者重新实现子模块，进而构建自己垂直方方向的爬虫。

本项目将爬虫的各个功能流程区分成Spider模块（主控），Downloader模块（下载器），Page_processer模块（页面分析），Scheduler模块（任务队列），Pipeline模块（结果输出）；


**执行过程简述**：
1. Spider从Scheduler中获取包含待抓取url的Request对象，传入Downloader，Downloader下载该url所对应的页面或者其他类型的数据，生成Page对象；
2. Spider调用PageProcesser模块解析Page中的数据，并存入Page中的PageItems中（以Key-Value对的形式保存），同时存入解析结果中的待抓取链接，Spider会将待抓取链接存入Scheduler模块中的Request队列中；
3. Spider调用Pipeline模块输出Page中的PageItems的结果;
4. 执行步骤1，直至所有链接被处理完成，则Spider被挂起等待下一个待抓取链接或者终止。



![image](https://github.com/hu17889/doc/blob/master/go_spider/img/project.png)


## 模块


### [Downloader](http://godoc.org/github.com/hu17889/go_spider/core/downloader)

**功能**：Spider从Scheduler中获取包含待抓取url的Request对象，传入Downloader，Downloader下载该url所对应的页面或者其他类型的数据，现在支持html和json两种结果类型或者无结果类型，生成Page对象，同时找到下载结果所对应的解析go包并生成解析器存入Page对象中，如html是[goquery包](https://github.com/PuerkitoBio/goquery)，json数据是[simplejson包](https://github.com/bitly/go-simplejson/blob/master/simplejson.go)。




