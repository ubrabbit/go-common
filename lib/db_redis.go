package lib

/*
http://godoc.org/github.com/garyburd/redigo/redis

封装常用的Redis命令
REDIS pool的使用
*/

import (
	"fmt"
	"time"
)

import (
	redis "github.com/garyburd/redigo/redis"

	. "github.com/ubrabbit/go-common/common"
)

var (
	g_RedisPool *RedisPool = nil
)

const (
	MAX_REDIS_POOL_ACTIVE = 4096
)

type poolDial struct {
	Conn redis.Conn
	Err  error
}

type RedisPool struct {
	host string
	port int
	conn redis.Conn
	pool *redis.Pool
}

func newRedisPool(host string) *redis.Pool {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			ch := make(chan poolDial)
			defer close(ch)
			go func() {
				for {
					//在短期高并发导致端口用尽时，会报 cannot assign requested address 或 too many open files 错误
					//所以需要用chan等待连接释放
					c, err := redis.Dial("tcp", host)
					if err != nil {
						LogError("Connect Redis Error: %v", err)
						time.Sleep(1 * time.Second)
						continue
					}
					ch <- poolDial{c, err}
					break
				}
			}()

			rlt := <-ch
			return rlt.Conn, rlt.Err
		},
		MaxIdle:     10,
		MaxActive:   MAX_REDIS_POOL_ACTIVE, // 最大连接数量，如果不设置这个值默认就是无限，当短时间高并发时报：too many open files
		Wait:        true,                  // 当达到最大连接数量时，阻塞， 如果不加这个参数，会报: connection pool exhausted
		IdleTimeout: 360 * time.Second,
		//获取连接对象前检查下连接是否还活着
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

func InitRedis(host string, port int) *RedisPool {
	address := fmt.Sprintf("%s:%d", host, port)
	LogInfo("Connect Redis %s", address)
	conn, err := redis.Dial("tcp", address)
	CheckFatal(err)

	LogInfo("Connect Redis Succ")
	g_RedisPool = &RedisPool{conn: conn, pool: newRedisPool(address)}
	g_RedisPool.host = host
	g_RedisPool.port = port
	return g_RedisPool
}

func (self *RedisPool) String() string {
	return fmt.Sprintf("[Redis] %s:%d", self.host, self.port)
}

func (self *RedisPool) Close() {
	defer func() {
		err := recover()
		if err != nil {
			LogError("Close Redis Error: %v", err)
		}
	}()
	self.conn.Close()
	self.pool.Close()
}

func (self *RedisPool) GetConn() (redis.Conn, error) {
	conn := self.pool.Get()
	if conn != nil {
		return conn, nil
	}
	LogError("Redis Pool Cant Get Conn")
	if self.conn == nil {
		return nil, fmt.Errorf("DB RedisConn %v is not inited!", self.conn)
	}
	_, err := self.conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("DB RedisConn %v is not alived!", self.conn)
	}
	return self.conn, nil
}

func (self *RedisPool) doCmd(cmd string, arg ...interface{}) (interface{}, error) {
	conn, err := self.GetConn()
	if err != nil {
		return nil, err
	}
	//不加这行语句会导致死锁
	//比如同一个函数执行了两次 RedisExec，但获取的是不同的conn的情况
	defer conn.Close()
	result, err := conn.Do(cmd, arg...)
	return result, err
}

func RedisExec(cmd string, arg ...interface{}) interface{} {
	result, err := g_RedisPool.doCmd(cmd, arg...)
	if err != nil {
		LogPanic("RedisExec Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return result
}

func RedisGetString(cmd string, arg ...interface{}) interface{} {
	value, err := redis.String(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetString Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return value
}

func RedisGetInt(cmd string, arg ...interface{}) interface{} {
	value, err := redis.Int64(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetInt Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return value
}

func RedisGetList(cmd string, arg ...interface{}) []string {
	value_list, err := redis.Values(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetList Error: %s %v %v", cmd, arg, err)
		return nil
	}
	result := make([]string, 0)
	for _, value := range value_list {
		result = append(result, string(value.([]byte)))
	}
	return result
}

func RedisGetStringMap(cmd string, arg ...interface{}) map[string]string {
	value, err := redis.StringMap(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetMap Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return value
}

func RedisGetIntMap(cmd string, arg ...interface{}) map[string]int {
	value, err := redis.IntMap(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetMap Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return value
}

func RedisGetInt64Map(cmd string, arg ...interface{}) map[string]int64 {
	value, err := redis.Int64Map(g_RedisPool.doCmd(cmd, arg...))
	if err != nil {
		LogPanic("RedisGetMap Error: %s %v %v", cmd, arg, err)
		return nil
	}
	return value
}
