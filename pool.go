package redisx

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Cacher struct {
	pool   *redis.Pool
	Prefix string
}

func (c *Cacher) closeIfDown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("close")
	c.pool.Close()
}

type ProxyConf struct {
	AddrList       []string
	MaxActive      int
	MaxIdle        int
	Downgrade      bool
	Network        string
	Password       string
	Db             int
	Prefix         string
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	Wait           bool
}

type SentinelConf struct {
	MasterName string
	Master     *ProxyConf
	Slave      *ProxyConf
	*ProxyConf
}

type Conf struct {
	Mode         string                // 模式，支持 default/sentinel
	ProxyConf    map[string]*ProxyConf // 地址配置
	SentinelConf *SentinelConf         // 哨兵配置
}

var RedisClusterMap = make(map[string]*Cacher)

func InitRedis(conf Conf) (err error) {
	if conf.Mode == "sentinel" {
		initSentinel(conf.SentinelConf)
	} else {
		initNormal(conf.ProxyConf)
	}
	return nil
}

func setDefaultOpts(opts *ProxyConf) {

	if opts.Network == "" {
		opts.Network = "tcp"
	}

	if opts.MaxIdle == 0 {
		opts.MaxIdle = 3
	}

	if opts.MaxActive == 0 {
		opts.MaxIdle = 8
	}

	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 10 * time.Second
	}

	if opts.ConnectTimeout == 0 {
		opts.ConnectTimeout = 10 * time.Second
	}

	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = 5 * time.Second
	}

	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = 5 * time.Second
	}
}

func CloseRedis() error {
	for _, v := range RedisClusterMap {
		v.pool.Close()
	}
	return nil
}

func GetPool(clusterList ...string) redis.Conn {

	clusterName := "default"
	if len(clusterList) > 0 {
		clusterName = clusterList[0]
	}

	cacher, ok := RedisClusterMap[clusterName]
	if !ok {
		cacher = RedisClusterMap["default"]
	}
	return cacher.pool.Get()
}

var NoAddr = errors.New("redisx:no addr")

func pool(opts *ProxyConf) func(addr string) *redis.Pool {
	return func(addr string) *redis.Pool {
		return poolWithDial(opts, dial(addr, opts))
	}
}

func poolWithDial(opts *ProxyConf, dial func() (redis.Conn, error)) *redis.Pool {
	return &redis.Pool{
		MaxIdle:      opts.MaxIdle,
		MaxActive:    opts.MaxActive,
		Wait:         opts.Wait,
		IdleTimeout:  opts.IdleTimeout,
		Dial:         dial,
		TestOnBorrow: testPing(),
	}
}

func dial(addr string, opts *ProxyConf) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {

		conn, err := redis.DialTimeout(opts.Network, addr, opts.ConnectTimeout, opts.ReadTimeout, opts.WriteTimeout)
		if err != nil {
			return nil, err
		}
		if opts.Password != "" {
			if _, err := conn.Do("AUTH", opts.Password); err != nil {
				conn.Close()
				return nil, err
			}
		}
		if opts.Db != 0 {
			if _, err := conn.Do("SELECT", opts.Db); err != nil {
				conn.Close()
				return nil, err
			}
		}
		return conn, err
	}
}

func testPing() func(conn redis.Conn, t time.Time) error {
	return func(conn redis.Conn, t time.Time) error {
		_, err := conn.Do("PING")
		return err
	}
}
