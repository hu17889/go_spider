// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package downloader

import (
    "core/common/page"
    "core/common/request"
)

type Downloader interface {
    Download(req *request.Request) *page.Page
}
