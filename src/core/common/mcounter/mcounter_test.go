package mcounter_test

import (
    "common/mcounter"
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
