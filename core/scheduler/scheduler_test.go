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
    r = request.NewRequest("http://baidu.com", "html")

    var s *scheduler.QueueScheduler
    s = scheduler.NewQueueScheduler()

    s.Push(r)
    var count int = s.Count()
    fmt.Println(count)

    var r1 *request.Request
    r1 = s.Poll()
    fmt.Println(r1)
}
