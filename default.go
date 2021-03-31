package redisx

import (
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

func initNormal(confs map[string]*ProxyConf) {
	for confName, conf := range confs {
		setDefaultOpts(conf)
		c := &Cacher{
			Prefix: conf.Prefix,
		}
		c.pool = getNormalPool(conf)
		go c.closeIfDown()
		RedisClusterMap[confName] = c
	}
}

func getNormalPool(opts *ProxyConf) *redis.Pool {
	setDefaultOpts(opts)
	return poolWithDial(opts, func() (redis.Conn, error) {
		addrs := opts.AddrList
		var addr string
		if len(addrs) == 0 {
			addr = "127.0.0.1:6379"
		} else {
			rand.Seed(time.Now().UnixNano())
			addr = addrs[rand.Intn(len(addrs))]
		}
		return dial(addr, opts)()
	})
}
