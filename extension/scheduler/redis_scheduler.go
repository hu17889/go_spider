package scheduler

import (
    "encoding/json"
    "github.com/garyburd/redigo/redis"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/request"
    "sync"
)

type RedisScheduler struct {
    locker      *sync.Mutex
    requestList string
    redisAddr   string
    redisPool   *redis.Pool
    maxConn     int
}

func NewRedisScheduler(addr string, maxConn int) *RedisScheduler {
    rs := &RedisScheduler{redisAddr: addr}
    rs = rs.Init()
    return rs
}

func (this *RedisScheduler) Init() *RedisScheduler {
    this.redisPool = redis.NewPool(this.newConn, this.maxConn)

    this.locker = new(sync.Mutex)
    this.requestList = "go_spider_request"
    return this
}

func (this *RedisScheduler) newConn() (redis.Conn, error) {
    return redis.Dial("tcp", this.redisAddr)
}
func (this *RedisScheduler) Push(requ *request.Request) {
    this.locker.Lock()
    defer this.locker.Unlock()

    requJson, err := json.Marshal(requ)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
        return
    }

    conn := this.redisPool.Get()
    //defer this.redisPool.Close()

    _, err = conn.Do("RPUSH", this.requestList, requJson)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
        return
    }
}

func (this *RedisScheduler) Poll() *request.Request {
    this.locker.Lock()
    defer this.locker.Unlock()

    conn := this.redisPool.Get()
    //defer this.redisPool.Close()

    length, err := conn.Do("LLEN", this.requestList)

    if length.(int64) <= 0 {
        return nil
    }
    buf, err := conn.Do("LPOP", this.requestList)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Poll Error: " + err.Error())
        return nil
    }

    requ := &request.Request{}

    err = json.Unmarshal(buf.([]byte), requ)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Poll Error: " + err.Error())
        return nil
    }

    return requ
}

func (this *RedisScheduler) Count() int {
    this.locker.Lock()

    conn := this.redisPool.Get()
    //defer this.redisPool.Close()
    len, err := conn.Do("LLEN", this.requestList)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Count Error: " + err.Error())
        this.locker.Unlock()
        return 0
    }
    this.locker.Unlock()
    return len.(int)
}
