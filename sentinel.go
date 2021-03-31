package redisx

import (
	"math/rand"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
)

func initSentinel(conf *SentinelConf) {

	sentinel := getSentinel(conf)
	master := &Cacher{
		Prefix: conf.Prefix,
	}
	master.pool = getSentinelMasterPool(sentinel, conf.Master)
	RedisClusterMap["default"] = master
	go master.closeIfDown()

	slave := &Cacher{
		Prefix: conf.Prefix,
	}
	slave.pool = getSentinelSlavePool(sentinel, conf.Slave)
	RedisClusterMap["read"] = slave
	go slave.closeIfDown()
}

func getSentinel(opts *SentinelConf) *sentinel.Sentinel {
	setDefaultOpts(opts.ProxyConf)
	if opts.MasterName == "" {
		opts.MasterName = "mymaster"
	}
	return &sentinel.Sentinel{
		Addrs:      opts.AddrList,
		MasterName: opts.MasterName,
		Pool:       pool(opts.ProxyConf),
	}
}

func getSentinelMasterPool(sntnl *sentinel.Sentinel, opts *ProxyConf) *redis.Pool {
	setDefaultOpts(opts)
	return poolWithDial(opts, func() (redis.Conn, error) {
		masterAddr, err := sntnl.MasterAddr()
		if err != nil {
			return nil, err
		}
		return dial(masterAddr, opts)()
	})
}

func getSentinelSlavePool(sntnl *sentinel.Sentinel, opts *ProxyConf) *redis.Pool {
	setDefaultOpts(opts)
	return poolWithDial(opts, func() (redis.Conn, error) {
		slaveAddrs, err := sntnl.SlaveAddrs()
		if err != nil {
			return nil, err
		}
		rand.Seed(time.Now().UnixNano())
		slaveAddr := slaveAddrs[rand.Intn(len(slaveAddrs))]
		return dial(slaveAddr, opts)()
	})
}
