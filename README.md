go_spider
=========

**v0.1** First Milestone，完成基本框架

## 简介


本项目基于golang开发，是一个垂直领域的爬虫引擎，主要希望能将各个功能模块区分开，方便使用者重新实现子模块，进而构建自己垂直方方向的爬虫。

本项目将爬虫的各个功能流程区分成Spider模块（主控），Downloader模块（下载器），Page_processer模块（页面分析），Scheduler模块（任务队列），Pipeline模块（结果输出）；


**调用过程**：Spider从Scheduler中获取包含待抓取url的Request对象，传入Downloader，Downloader下载该url所对应的页面或者其他类型的数据，现在支持html和json两种类型，生成Page对象



![image](https://github.com/hu17889/doc/blob/master/go_spider/img/project.png)


## 模块


### Downloader

**功能**：Spider从Scheduler中获取包含待抓取url的Request对象，传入Downloader，Downloader下载该url所对应的页面或者其他类型的数据，现在支持html和json两种类型，生成Page对象，同时找到下载结果所对应的解析go包并生成解析器存入Page中，如html是[goquery包](https://github.com/PuerkitoBio/goquery)，json数据是[simplejson包](https://github.com/bitly/go-simplejson/blob/master/simplejson.go)
