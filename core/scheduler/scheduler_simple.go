// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The package is useless
package scheduler

import (
    "github.com/hu17889/go_spider/core/common/request"
)

type SimpleScheduler struct {
    queue chan *request.Request
}

func NewSimpleScheduler() *SimpleScheduler {
    ch := make(chan *request.Request, 1024)
    return &SimpleScheduler{ch}
}

func (this *SimpleScheduler) Push(requ *request.Request) {
    this.queue <- requ
}

func (this *SimpleScheduler) Poll() *request.Request {
    if len(this.queue) == 0 {
        return nil
    } else {
        return <-this.queue
    }
}

func (this *SimpleScheduler) Count() int {
    return len(this.queue)
}
