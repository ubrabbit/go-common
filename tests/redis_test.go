package tests

import (
	"fmt"
	"testing"
)

import (
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
)

func TestRedis(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRedis  =====================\n")

	config.InitConfig("config_test.conf")
	cfg := config.GetRedisConfig()
	pool := lib.InitRedis(cfg.IP, cfg.Port)

	lib.RedisExec("del", "Score")
	lib.RedisExec("del", "Score_NotExists")
	lib.RedisExec("del", "Name")
	lib.RedisExec("del", "LIST")

	lib.RedisExec("hset", "Score", "test", "1")
	lib.RedisExec("set", "Name", "Hello")
	lib.RedisExec("rpush", "LIST", "1")

	rlt := lib.RedisExec("hset", "Score", "test2", "2")
	fmt.Println("RedisExec : ", rlt)
	fmt.Println("RedisGetString : ", lib.RedisGetString("get", "Name"))
	fmt.Println("RedisGetInt : ", lib.RedisGetInt("hget", "Score", "test"))
	fmt.Println("RedisGetList1 : ", lib.RedisGetList("lrange", "LIST", 0, -1))
	fmt.Println("RedisGetList2 : ", lib.RedisGetList("hgetall", "Score"))
	fmt.Println("RedisGetIntMap : ", lib.RedisGetIntMap("hgetall", "Score"))

	pool.Close()
}
