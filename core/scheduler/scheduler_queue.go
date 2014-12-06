// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package scheduler

import (
    "container/list"
    "crypto/md5"
    "github.com/hu17889/go_spider/core/common/request"
    "sync"
    //"fmt"
)

type QueueScheduler struct {
    locker *sync.Mutex
    rm     bool
    rmKey  map[[md5.Size]byte]*list.Element
    queue  *list.List
}

func NewQueueScheduler(rmDuplicate bool) *QueueScheduler {
    queue := list.New()
    rmKey := make(map[[md5.Size]byte]*list.Element)
    locker := new(sync.Mutex)
    return &QueueScheduler{rm: rmDuplicate, queue: queue, rmKey: rmKey, locker: locker}
}

func (this *QueueScheduler) Push(requ *request.Request) {
    this.locker.Lock()
    var key [md5.Size]byte
    if this.rm {
        key = md5.Sum([]byte(requ.GetUrl()))
        if _, ok := this.rmKey[key]; ok {
            this.locker.Unlock()
            return
        }
    }
    e := this.queue.PushBack(requ)
    if this.rm {
        this.rmKey[key] = e
    }
    this.locker.Unlock()
}

func (this *QueueScheduler) Poll() *request.Request {
    this.locker.Lock()
    if this.queue.Len() <= 0 {
        this.locker.Unlock()
        return nil
    }
    e := this.queue.Front()
    requ := e.Value.(*request.Request)
    key := md5.Sum([]byte(requ.GetUrl()))
    this.queue.Remove(e)
    if this.rm {
        delete(this.rmKey, key)
    }
    this.locker.Unlock()
    return requ
}

func (this *QueueScheduler) Count() int {
    this.locker.Lock()
    len := this.queue.Len()
    this.locker.Unlock()
    return len
}
