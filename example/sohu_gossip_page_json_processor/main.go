// example for request extension

package main

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/hu17889/go_spider/core/common/page"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/spider"

    "encoding/json"
    "fmt"
    "regexp"
    "strings"
)

const (
    wkSohuUrl       = "http://yule.sohu.com/gossip/index.shtml"
    wkSohuYule      = `http://changyan.sohu.com/node/html?appid=cyqemw6s1&client_id=cyqemw6s1&topicsid=%s&spSize=5`
    wkSohuPic       = `http://changyan.sohu.com/node/html?appid=cyqemw6s1&client_id=cyqemw6s1&topicsid=9000%s&spSize=5`
    maxWKSouhuLayer = 3 // 最多抓取多少页
)

var rxYule = regexp.MustCompile(`^http://yule\.sohu\.com/.*?/n(.*?).shtml`)      // 筛选出娱乐版块
var rxPic = regexp.MustCompile(`^http://pic\.yule\.sohu\.com/group-(.*?).shtml`) // 筛选出图片版块；八卦版块混合了这两个

type MyPageProcesser struct {
}

type ChangyanListDataJson struct {
    OuterCmtSum      int `json:"outer_cmt_sum"`
    ParticipationSum int `json:"participation_sum"`
}

type ChangyanJson struct {
    ListData ChangyanListDataJson `json:"listData"`
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

func addRequest(p *page.Page, tag, url, cookie, content string) {
    req := request.NewRequest(url, "json", tag, "GET", "", nil, nil, nil, content)
    p.AddTargetRequestWithParams(req)
}

func (this MyPageProcesser) Process(p *page.Page) {
    query := p.GetHtmlParser()

    if p.GetUrlTag() == "index" {
        query.Find(`div[class="main area"] div[class="lc"] ul li a`).Each(func(i int, s *goquery.Selection) {
            url, isExsit := s.Attr("href")
            if isExsit {
                reg := regexp.MustCompile(`^do not know what is this`)
                var fmtStr string
                if rxYule.MatchString(url) {
                    reg = rxYule
                    fmtStr = wkSohuYule
                }

                if rxPic.MatchString(url) {
                    reg = rxPic
                    fmtStr = wkSohuPic
                }

                regxpArrag := reg.FindStringSubmatch(url)
                if len(regxpArrag) == 2 {
                    addRequest(p, "changyan", fmt.Sprintf(fmtStr, regxpArrag[1]), "", s.Text())
                }
            }
        })
    }

    if p.GetUrlTag() == "changyan" {
        jsonMap := ChangyanJson{}
        err := json.NewDecoder(strings.NewReader(p.GetBodyStr())).Decode(&jsonMap)
        if err == nil {
            content, ok := p.GetRequest().GetExtension().(string)
            if ok {
                fmt.Println("标题:", content, " 评论数:", jsonMap.ListData.OuterCmtSum, " 参与数:", jsonMap.ListData.ParticipationSum)
            }
        }
    }
}

func main() {
    req := request.NewRequest(wkSohuUrl, "html", "index", "GET", "", nil, nil, nil, nil)
    sohuSpider := spider.NewSpider(NewMyPageProcesser(), "Sohu").
        AddRequest(req).
        SetSleepTime("rand", 500, 1000).
        SetThreadnum(2)

    for i := 1; i < maxWKSouhuLayer; i++ {
        url := fmt.Sprintf("http://yule.sohu.com/gossip/index_%d.shtml", 5301-i) // 一个神奇的数字，
        req := request.NewRequest(url, "html", "index", "GET", "", nil, nil, nil, nil)
        sohuSpider.AddRequest(req)
    }

    sohuSpider.Run()
}
