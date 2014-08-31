//
package resource_manage

type ResourceManage interface {
    GetOne()
    FreeOne()
    Has() int
    Left() int
}
