package redisx_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	. "github.com/uccu/go-redisx"
)

func Test1(t *testing.T) {

	conf := Conf{
		Mode: "default",
		ProxyConf: map[string]*ProxyConf{
			"default": {
				AddrList: []string{"127.0.0.1:6379"},
			},
		},
	}

	InitRedis(conf)

	conn := GetPool()
	conn.Send("SET", redis.Args{}.Add("testKey1").AddFlat(1)...)
	reply, err := redis.String(conn.Do("GET", redis.Args{}.Add("testKey1")...))
	conn.Close()
	fmt.Println(reply, err)

}
