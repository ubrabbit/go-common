package tests

import (
	"fmt"
	"testing"
	"time"
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

	lib.RedisExec("hset", "Score", "test", "10086")
	lib.RedisExec("set", "Name", "Hello")
	lib.RedisExec("rpush", "LIST", "123")

	rlt := lib.RedisExec("hset", "Score", "aaa", "123")
	fmt.Println("RedisExec : ", rlt)
	fmt.Println("RedisGetString : ", lib.RedisGetString("get", "Name"))
	fmt.Println("RedisGetInt : ", lib.RedisGetInt("hget", "Score", "aaa"))
	fmt.Println("RedisGetList1 : ", lib.RedisGetList("lrange", "LIST", 0, -1))
	fmt.Println("RedisGetList2 : ", lib.RedisGetList("hgetall", "Score"))

	pool.Close()
	lib.RedisExec("set", "Name", "Hello")
	fmt.Println("RedisGetList2 After Closed: ", lib.RedisGetList("hgetall", "Score"))
	time.Sleep(1 * time.Second)
}
