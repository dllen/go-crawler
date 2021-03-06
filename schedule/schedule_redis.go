package schedule

import (
	"time"

	"github.com/dllen/go-crawler/common"
	"github.com/dllen/go-crawler/config"
	"github.com/dllen/go-crawler/logger"
	"github.com/dllen/go-crawler/model"
	"github.com/garyburd/redigo/redis"
)

type RedisSchedule struct {
	name    string
	address string
	pool    *redis.Pool
}

func NewRedisSchedule(config *config.Config) Schedule {
	schedule := &RedisSchedule{}
	schedule.address = config.RedisAddr
	schedule.name = config.Name
	schedule.connect()

	return schedule
}

func (r *RedisSchedule) connect() {
	r.pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", r.address) },
	}

	go r.CronCount(1)
}

func (r *RedisSchedule) Push(req *model.Request) {
	conn := r.pool.Get()
	defer conn.Close()

	ParseReqs := common.ParseReq([]*model.Request{req}, nil)
	for _, req := range ParseReqs {
		logger.Info("Push URL:", req.Url, req.ProcessName)
		body, err := req.Write()
		if err != nil {
			logger.Info("Push URL:", err.Error())
			continue
		}
		_, err = conn.Do("LPUSH", r.name, body)
		if err != nil {
			logger.Info("Push URL:", err.Error())
			continue
		}
	}
}

func (r *RedisSchedule) PushMulti(reqs []*model.Request) {
	conn := r.pool.Get()
	defer conn.Close()

	ParseReqs := common.ParseReq(reqs, nil)
	for _, req := range ParseReqs {
		logger.Info("Push URL:", req.Url, req.ProcessName)
		body, err := req.Write()
		if err != nil {
			logger.Info("Push URL:", err.Error())
			continue
		}
		_, err = conn.Do("LPUSH", r.name, body)
		if err != nil {
			logger.Info("Push URL:", err.Error())
			continue
		}
	}
}

func (r *RedisSchedule) Pop() (*model.Request, bool) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err := redis.ByteSlices(conn.Do("BRPOP", r.name, 5))
	if err != nil {
		logger.Info("Pop Url: ", err.Error())
		return nil, true
	}

	req := &model.Request{}
	if err := req.Read(value[1]); err != nil {
		logger.Info("Pop Url: ", err.Error())
		return nil, true
	}

	logger.Info("Pop Url:", req.Url, req.ProcessName)
	return req, true
}

func (r *RedisSchedule) Count() int {
	conn := r.pool.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("LLEN", r.name))
	if err != nil {
		logger.Info("Count  ", err.Error())
		return -1
	}
	return value
}

func (r *RedisSchedule) Close() {
	r.pool.Close()
}

func (r *RedisSchedule) CronCount(flushTime int) {
	ticker := time.NewTicker(time.Second * time.Duration(flushTime))
	go func() {
		for range ticker.C {
			logger.Info("RedisSchedule Count:", r.Count())
		}
	}()
}

func init() {
	RegisterSchedule("redis", NewRedisSchedule)
}
