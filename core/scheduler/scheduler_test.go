// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package scheduler_test

import (
    "fmt"
    "github.com/hu17889/go_spider/core/common/request"
    "github.com/hu17889/go_spider/core/scheduler"
    "testing"
)

func TestQueueScheduler(t *testing.T) {
    var r *request.Request
    r = request.NewRequest("http://baidu.com", "html", "", "GET", "", nil, nil, nil, nil)
    fmt.Printf("%v\n", r)

    var s *scheduler.QueueScheduler
    s = scheduler.NewQueueScheduler(false)

    s.Push(r)
    var count int = s.Count()
    if count != 1 {
        t.Error("count error")
    }
    fmt.Println(count)

    var r1 *request.Request
    r1 = s.Poll()
    if r1 == nil {
        t.Error("poll error")
    }
    fmt.Printf("%v\n", r1)

    // remove duplicate
    s = scheduler.NewQueueScheduler(true)

    r2 := request.NewRequest("http://qq.com", "html", "", "GET", "", nil, nil, nil, nil)
    s.Push(r)
    s.Push(r2)
    s.Push(r)
    count = s.Count()
    if count != 2 {
        t.Error("count error")
    }
    fmt.Println(count)

    r1 = s.Poll()
    if r1 == nil {
        t.Error("poll error")
    }
    fmt.Printf("%v\n", r1)
    r1 = s.Poll()
    if r1 == nil {
        t.Error("poll error")
    }
    fmt.Printf("%v\n", r1)
}
