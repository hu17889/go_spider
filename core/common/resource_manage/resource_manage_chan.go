package resource_manage

// ResourceManageChan inherits the ResourceManage interface.
// In spider, ResourceManageChan manage resource of Coroutine to crawl page.
type ResourceManageChan struct {
    capnum int
    mc     chan int
}

// NewResourceManageChan returns initialized ResourceManageChan object which contains a resource pool.
// The num is the resource limit.
func NewResourceManageChan(num int) *ResourceManageChan {
    mc := make(chan int, num)
    return &ResourceManageChan{mc: mc, capnum: num}
}

// The GetOne apply for one resource.
// If resource pool is empty, current coroutine will be blocked.
func (this *ResourceManageChan) GetOne() {
    this.mc <- 1
}

// The FreeOne free resource and return it to resource pool.
func (this *ResourceManageChan) FreeOne() {
    <-this.mc
}

// The Has query for how many resource has been used.
func (this *ResourceManageChan) Has() int {
    return len(this.mc)
}

// The Left query for how many resource left in the pool.
func (this *ResourceManageChan) Left() int {
    return this.capnum - len(this.mc)
}
