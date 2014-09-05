// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package scheduler

import (
    "github.com/hu17889/go_spider/core/common/request"
)

type Scheduler interface {
    Push(requ *request.Request)
    Poll() *request.Request
    Count() int
}
