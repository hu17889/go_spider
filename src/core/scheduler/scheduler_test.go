package scheduler_test

import (
    "common/request"
    "fmt"
    "scheduler"
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
