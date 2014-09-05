//
package resource_manage_test

import (
    "core/common/resource_manage"
    "testing"
)

func TestResourceManage(t *testing.T) {
    var mc *mcounter.ResourceManage
    mc = mcounter.NewResourceManageChan(1)
    mc.GetOne()
    println("incr")
    mc.FreeOne()
    println("decr")
    mc.GetOne()
    println("incr")
}
