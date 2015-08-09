package scheduler

import (
    "encoding/json"
    "github.com/garyburd/redigo/redis"
    "github.com/hu17889/go_spider/core/common/mlog"
    "github.com/hu17889/go_spider/core/common/request"
    "sync"
)

type RedisScheduler struct {
    locker                *sync.Mutex
    requestList           string
    urlList               string
    redisAddr             string
    redisPool             *redis.Pool
    maxConn               int
    maxIdle               int
    forbiddenDuplicateUrl bool
    queueMax              int
}

func NewRedisScheduler(addr string, maxConn, maxIdle int, forbiddenDuplicateUrl bool) *RedisScheduler {
    rs := &RedisScheduler{
        redisAddr:             addr,
        forbiddenDuplicateUrl: forbiddenDuplicateUrl,
        maxConn:               maxConn,
        maxIdle:               maxIdle,
        requestList:           "go_spider_request",
        urlList:               "go_spider_url",
    }
    rs = rs.Init()
    return rs
}

func (this *RedisScheduler) Init() *RedisScheduler {
    this.redisPool = redis.NewPool(this.newConn, this.maxIdle)
    this.redisPool.MaxActive = this.maxConn
    this.locker = new(sync.Mutex)
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
    defer conn.Close()

    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
        return
    }
    if this.forbiddenDuplicateUrl {
        urlExist, err := conn.Do("HGET", this.urlList, requ.GetUrl())
        if err != nil {
            mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
            return
        }
        if urlExist != nil {
            return
        }

        conn.Do("MULTI")
        _, err = conn.Do("HSET", this.urlList, requ.GetUrl(), 1)
        if err != nil {
            mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
            conn.Do("DISCARD")
            return
        }
    }
    _, err = conn.Do("RPUSH", this.requestList, requJson)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
        if this.forbiddenDuplicateUrl {
            conn.Do("DISCARD")
        }
        return
    }

    if this.forbiddenDuplicateUrl {
        conn.Do("EXEC")
    }
}

func (this *RedisScheduler) Poll() *request.Request {
    this.locker.Lock()
    defer this.locker.Unlock()

    conn := this.redisPool.Get()
    defer conn.Close()

    length, err := this.count()
    if err != nil {
        return nil
    }
    if length <= 0 {
        mlog.LogInst().LogError("RedisScheduler Poll length 0")
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
    defer this.locker.Unlock()
    var length int
    var err error

    length, err = this.count()
    if err != nil {
        return 0
    }

    return length
}

func (this *RedisScheduler) count() (int, error) {
    conn := this.redisPool.Get()
    defer conn.Close()
    length, err := conn.Do("LLEN", this.requestList)
    if err != nil {
        mlog.LogInst().LogError("RedisScheduler Count Error: " + err.Error())
        return 0, err
    }
    return int(length.(int64)), nil
}
