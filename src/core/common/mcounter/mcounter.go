package mcounter

type Mcounter interface {
    Incr()
    Decr()
    Count() int
}
