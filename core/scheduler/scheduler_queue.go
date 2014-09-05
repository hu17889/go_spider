// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package scheduler

import (
    "github.com/hu17889/go_spider/core/common/request"
)

type QueueScheduler struct {
    queue chan *request.Request
}

func NewQueueScheduler() *QueueScheduler {
    ch := make(chan *request.Request, 1024)
    return &QueueScheduler{ch}
}

func (this *QueueScheduler) Push(requ *request.Request) {
    this.queue <- requ
}

func (this *QueueScheduler) Poll() *request.Request {
    if len(this.queue) == 0 {
        return nil
    } else {
        return <-this.queue
    }
}

func (this *QueueScheduler) Count() int {
    return len(this.queue)
}
