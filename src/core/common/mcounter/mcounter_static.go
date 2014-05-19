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
