// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package mcounter_test

import (
    "core/common/mcounter"
    "testing"
)

func TestMcounter(t *testing.T) {
    var mc *mcounter.StaticMcounter
    mc = mcounter.NewStaticMcounter(1)
    mc.Incr()
    println("incr")
    mc.Decr()
    println("decr")
    mc.Incr()
    println("incr")

    var mc1 mcounter.Mcounter
    mc1 = mcounter.NewStaticMcounter(1)
    mc1.Incr()
    println("incr")
    mc1.Decr()
    println("decr")
    mc1.Incr()
    println("incr")
}
