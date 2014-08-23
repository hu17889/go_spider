// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package mcounter

type StaticMcounter struct {
    mc chan int
}

func NewStaticMcounter(num int) *StaticMcounter {
    mc := make(chan int, num)
    return &StaticMcounter{mc}
}

func (this *StaticMcounter) Incr() {
    this.mc <- 1
}

func (this *StaticMcounter) Decr() {
    <-this.mc
}

func (this *StaticMcounter) Count() int {
    return len(this.mc)
}
