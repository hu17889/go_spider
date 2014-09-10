// Package resource_manage implements a resource management.
package resource_manage

// ResourceManage is an interface that who want implement an management object can realize these functions.
type ResourceManage interface {
    GetOne()
    FreeOne()
    Has() uint
    Left() uint
}
