// Package downloader is the main module of GO_SPIDER for download page.
package downloader

import (
    "core/common/page"
    "core/common/request"
)

type Downloader interface {
    Download(req *request.Request) *page.Page
}
