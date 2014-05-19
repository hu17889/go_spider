package downloader

import (
    "common/page"
    "common/request"
)

type Downloader interface {
    Download(req *request.Request) *page.Page
}
