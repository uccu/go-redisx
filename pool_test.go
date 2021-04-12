package redisx_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	. "github.com/uccu/go-redisx"
)

// ceshi
func Test1(t *testing.T) {

	pools := InitRedis([]*Conf{
		{
			Mode: "single",
			SingleConf: &ProxyConf{
				AddrList: []string{"127.0.0.1:6379"},
				Db:       0,
			},
		},
		{
			Mode: "single",
			SingleConf: &ProxyConf{
				Name:     "s1",
				AddrList: []string{"127.0.0.1:6379"},
				Db:       1,
			},
		},
	})
	conn := pools.GetPool()
	conn.Send("SET", redis.Args{}.Add("testKey1").AddFlat(1)...)

	conn2 := pools.GetPool("s1")
	conn2.Send("SET", redis.Args{}.Add("testKey1").AddFlat(2)...)

	reply, err := redis.String(conn.Do("GET", redis.Args{}.Add("testKey1")...))
	reply2, err := redis.String(conn2.Do("GET", redis.Args{}.Add("testKey1")...))

	conn.Close()
	fmt.Println(reply, reply2, err)

}
