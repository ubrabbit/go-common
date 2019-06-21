package tests

//go test -bench=".*" -count=1

import (
	"fmt"
	"testing"
)

import (
	lib "github.com/ubrabbit/go-public/lib"
)

const Redis_RunTimes = 100

func Benchmark_RedisSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			b.StopTimer()
			name := fmt.Sprintf("Name_%d", j+1)
			value := fmt.Sprintf("Value_%d", j+1)
			b.StartTimer()
			lib.RedisExec("set", name, value)
		}
	}
}

func Benchmark_RedisHSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			b.StopTimer()
			name := fmt.Sprintf("Name_%d", j+1)
			value := fmt.Sprintf("Value_%d", j+1)
			b.StartTimer()
			lib.RedisExec("hset", "Score", name, value)
		}
	}
}

func Benchmark_RedisRPush(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			lib.RedisExec("rpush", "LIST", j+1)
		}
	}
}

func Benchmark_RedisGetString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			b.StopTimer()
			name := fmt.Sprintf("Name_%d", j+1)
			b.StartTimer()
			lib.RedisGetString("get", name)
		}
	}
}

func Benchmark_RedisGetList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			lib.RedisGetList("lrange", "LIST", 0, -1)
		}
	}
}

func Benchmark_RedisGetList2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			lib.RedisGetList("hgetall", "Score")
		}
	}
}

func Benchmark_RedisGetMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < Redis_RunTimes; j++ {
			lib.RedisGetStringMap("hgetall", "Score")
		}
	}
}

func Benchmark_RedisPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 512; j++ {
			b.StopTimer()
			name := fmt.Sprintf("Name_%d", j+1)
			b.StartTimer()
			go lib.RedisGetString("get", name)
		}
	}
}

func init() {
	lib.InitRedis("127.0.0.1", 6379)
	lib.RedisExec("del", "Score")
	lib.RedisExec("del", "Name")
	lib.RedisExec("del", "LIST")
}
