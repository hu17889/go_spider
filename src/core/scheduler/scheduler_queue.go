package scheduler

import (
    "common/request"
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
