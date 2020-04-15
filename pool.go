package redis_timeseries_go

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type ConnPool interface {
	Get() redis.Conn
}

type SingleHostPool struct {
	*redis.Pool
}

func NewSingleHostPool(host string, authPass *string) *SingleHostPool {
	ret := &redis.Pool{
		Dial:         dialFuncWrapper(host, authPass),
		TestOnBorrow: testOnBorrow,
		MaxIdle:      maxConns,
	}

	return &SingleHostPool{ret}
}

type MultiHostPool struct {
	sync.Mutex
	pools    map[string]*redis.Pool
	hosts    []string
	authPass *string
}

func NewMultiHostPool(hosts []string, authPass *string) *MultiHostPool {
	return &MultiHostPool{
		pools:    make(map[string]*redis.Pool, len(hosts)),
		hosts:    hosts,
		authPass: authPass,
	}
}

func (p *MultiHostPool) Get() redis.Conn {
	p.Lock()
	defer p.Unlock()

	host := p.hosts[rand.Intn(len(p.hosts))]
	pool, found := p.pools[host]

	if !found {
		pool = &redis.Pool{
			Dial:         dialFuncWrapper(host, p.authPass),
			TestOnBorrow: testOnBorrow,
			MaxIdle:      maxConns,
		}
		p.pools[host] = pool
	}

	return pool.Get()
}

func dialFuncWrapper(host string, authPass *string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", host)
		if err != nil {
			return conn, err
		}
		if authPass != nil {
			_, err = conn.Do("AUTH", *authPass)
		}
		return conn, err
	}
}

func testOnBorrow(c redis.Conn, t time.Time) (err error) {
	if time.Since(t) > time.Millisecond {
		_, err = c.Do("PING")
	}
	return err
}
