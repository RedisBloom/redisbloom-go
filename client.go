package redis_timeseries_go

import (
	"strings"
	"github.com/gomodule/redigo/redis"
)

// TODO: refactor this hard limit and revise client locking
// Client Max Connections
var maxConns = 500

// Client is an interface to time series redis commands
type Client struct {
	Pool ConnPool
	Name string
}

// NewClient creates a new client connecting to the redis host, and using the given name as key prefix.
// Addr can be a single host:port pair, or a comma separated list of host:port,host:port...
// In the case of multiple hosts we create a multi-pool and select connections at random
func NewClient(addr, name string, authPass *string) *Client {
	addrs := strings.Split(addr, ",")
	var pool ConnPool
	if len(addrs) == 1 {
		pool = NewSingleHostPool(addrs[0], authPass)
	} else {
		pool = NewMultiHostPool(addrs, authPass)
	}
	ret := &Client{
		Pool: pool,
		Name: name,
	}
	return ret
}

// Add - Add (or create and add) a new value to the filter
// args:
// key - the name of the filter 
// item - the item to add 
func (client *Client) Add(key string, item string) (exists bool, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("BF.ADD", key, item))
}

// Exists - Determines whether an item may exist in the Bloom Filter or not. 
// args:
// key - the name of the filter 
// item - the item to check for 
func (client *Client) Exists(key string, item string) (exists bool, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("BF.EXISTS", key, item))
}
