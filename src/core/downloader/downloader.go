//
package downloader

import (
    "core/common/page"
    "core/common/request"
)

type Downloader interface {
    Download(req *request.Request) *page.Page
}
