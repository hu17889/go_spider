//
package resource_manage

type ResourceManageChan struct {
    capnum int
    mc     chan int
}

func NewResourceManageChan(num int) *ResourceManageChan {
    mc := make(chan int, num)
    return &ResourceManageChan{mc: mc, capnum: num}
}

func (this *ResourceManageChan) GetOne() {
    this.mc <- 1
}

func (this *ResourceManageChan) FreeOne() {
    <-this.mc
}

func (this *ResourceManageChan) Has() int {
    return len(this.mc)
}

func (this *ResourceManageChan) Left() int {
    return this.capnum - len(this.mc)
}
