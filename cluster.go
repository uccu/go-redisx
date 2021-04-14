package redisx

import (
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
)

type ClusterConf struct {
	*ProxyConf
}

var RefreshFailed = errors.New("redisx:refresh failed")

func initCluster(options *ClusterConf) Pool {

	cluster := &redisc.Cluster{
		StartupNodes: options.AddrList,
		CreatePool: func(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
			return poolWithDial(options.ProxyConf, dial(addr, options.ProxyConf)), nil
		},
	}

	if err := cluster.Refresh(); err != nil {
		panic(RefreshFailed)
	}

	return cluster

}
