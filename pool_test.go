package redisx_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	. "github.com/uccu/go-redisx"
)

// ceshi
func Test1(t *testing.T) {

	conf := Conf{
		Mode: "single",
		SingleConf: &ProxyConf{
			AddrList: []string{"127.0.0.1:6379"},
		},
	}

	pools := InitRedis(conf)
	conn := pools.GetPool()
	conn.Send("SET", redis.Args{}.Add("testKey1").AddFlat(2)...)
	reply, err := redis.String(conn.Do("GET", redis.Args{}.Add("testKey1")...))
	conn.Close()
	fmt.Println(reply, err)

}
