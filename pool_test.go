package redisx_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	. "github.com/uccu/go-redisx"
)

// ceshi
func TestSingle(t *testing.T) {

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

	reply, err1 := redis.String(conn.Do("GET", redis.Args{}.Add("testKey1")...))
	reply2, err2 := redis.String(conn2.Do("GET", redis.Args{}.Add("testKey1")...))

	conn.Close()
	fmt.Println(reply, reply2, err1, err2)

}

func TestCluster(t *testing.T) {

	pools := InitRedis([]*Conf{
		{
			Mode: "cluster",
			ClusterConf: &ClusterConf{
				ProxyConf{
					AddrList: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003"},
					Db:       0,
				},
			},
		},
	})

	conn := pools.GetPool()
	conn.Send("SET", redis.Args{}.Add("a").AddFlat(1)...)
	reply, err1 := redis.Int(conn.Do("GET", redis.Args{}.Add("a")...))

	conn2 := pools.GetPool()
	conn2.Send("SET", redis.Args{}.Add("b").AddFlat(1)...)
	reply2, err2 := redis.Int(conn2.Do("GET", redis.Args{}.Add("b")...))

	conn.Close()
	fmt.Println(reply, reply2, err1, err2)

}
