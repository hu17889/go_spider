// Package resource_manage implements a resource management.
package resource_manage


type ResourceManage interface {
    GetOne()
    FreeOne()
    Has() int
    Left() int
}
