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
	conn        redis.Conn
	requestList string
}

func NewRedisScheduler(conn redis.Conn) *RedisScheduler {
	locker := new(sync.Mutex)
	return &RedisScheduler{conn: conn, locker: locker, requestList: "go_spider_request"}
}

func (this *RedisScheduler) Push(requ *request.Request) {
	this.locker.Lock()
	defer this.locker.Unlock()

	requJson, err := json.Marshal(requ)
	if err != nil {
		mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
		return
	}

	_, err = this.conn.Do("RPUSH", this.requestList, requJson)
	if err != nil {
		mlog.LogInst().LogError("RedisScheduler Push Error: " + err.Error())
		return
	}
}

func (this *RedisScheduler) Poll() *request.Request {
	this.locker.Lock()
	defer this.locker.Unlock()
	length, err := this.conn.Do("LLEN", this.requestList)
	if length.(int64) <= 0 {
		return nil
	}
	buf, err := this.conn.Do("LPOP", this.requestList)
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
	len, err := this.conn.Do("LLEN", this.requestList)
	if err != nil {
		mlog.LogInst().LogError("RedisScheduler Count Error: " + err.Error())
		this.locker.Unlock()
		return 0
	}
	this.locker.Unlock()
	return len.(int)
}
