package redis_bloom_go

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
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

// NewClientFromPool creates a new Client with the given pool and client name
func NewClientFromPool(pool *redis.Pool, name string) *Client {
	ret := &Client{
		Pool: pool,
		Name: name,
	}
	return ret
}

// Reserve - Creates an empty Bloom Filter with a given desired error ratio and initial capacity.
// args:
// key - the name of the filter
// error_rate - the desired probability for false positives
// capacity - the number of entries you intend to add to the filter
func (client *Client) Reserve(key string, error_rate float64, capacity uint64) (err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("BF.RESERVE", key, strconv.FormatFloat(error_rate, 'g', 16, 64), capacity)
	return err
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

// Info - Return information about key
// args:
// key - the name of the filter
func (client *Client) Info(key string) (info map[string]int64, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	result, err := conn.Do("BF.INFO", key)
	if err != nil {
		return nil, err
	}

	values, err := redis.Values(result, nil)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("Info expects even number of values result")
	}
	info = map[string]int64{}
	for i := 0; i < len(values); i += 2 {
		key, err = redis.String(values[i], nil)
		if err != nil {
			return nil, err
		}
		info[key], err = redis.Int64(values[i+1], nil)
		if err != nil {
			return nil, err
		}
	}
	return info, nil
}
