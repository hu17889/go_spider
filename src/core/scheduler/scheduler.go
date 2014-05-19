package scheduler

import (
    "common/request"
)

type Scheduler interface {
    Push(requ *request.Request)
    Poll() *request.Request
    Count() int
}
