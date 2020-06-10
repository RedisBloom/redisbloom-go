package redis_bloom_go

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func getTestConnectionDetails() (string, string) {
	value, exists := os.LookupEnv("REDISBLOOM_TEST_HOST")
	host := "localhost:6379"
	password := ""
	valuePassword, existsPassword := os.LookupEnv("REDISBLOOM_TEST_PASSWORD")
	if exists && value != "" {
		host = value
	}
	if existsPassword && valuePassword != "" {
		password = valuePassword
	}
	return host, password
}

func createClient() *Client {
	host, password := getTestConnectionDetails()
	var ptr *string = nil
	if len(password) > 0 {
		ptr = &password
	}
	return NewClient(host, "test_client", ptr)
}

func TestNewClientFromPool(t *testing.T) {
	host, password := getTestConnectionDetails()
	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", host, redis.DialPassword(password))
	}, MaxIdle: maxConns}
	client1 := NewClientFromPool(pool, "bloom-client-1")
	client2 := NewClientFromPool(pool, "bloom-client-2")
	assert.Equal(t, client1.Pool, client2.Pool)
	err1 := client1.Pool.Close()
	err2 := client2.Pool.Close()
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

var client = createClient()
var _ = client.FlushAll()

var defaultDuration, _ = time.ParseDuration("1h")
var tooShortDuration, _ = time.ParseDuration("10ms")

func (client *Client) FlushAll() (err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("FLUSHALL")
	return err
}

func TestReserve(t *testing.T) {
	client.FlushAll()
	key := "test_RESERVE"
	err := client.Reserve(key, 0.1, 1000)
	assert.Nil(t, err)

	info, err := client.Info(key)
	assert.Nil(t, err)
	assert.Equal(t, info, map[string]int64{
		"Capacity":                 1000,
		"Expansion rate":           2,
		"Number of filters":        1,
		"Number of items inserted": 0,
		"Size":                     932,
	})

	err = client.Reserve(key, 0.1, 1000)
	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	client.FlushAll()
	key := "test_ADD"
	value := "test_ADD_value"
	exists, err := client.Add(key, value)
	assert.Nil(t, err)
	assert.True(t, exists)

	info, err := client.Info(key)
	assert.Nil(t, err)
	assert.NotNil(t, info)

	exists, err = client.Add(key, value)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestExists(t *testing.T) {
	client.FlushAll()
	client.Add("test_ADD", "test_EXISTS")

	exists, err := client.Exists("test_ADD", "test_EXISTS")
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = client.Exists("test_ADD", "test_EXISTS1")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestClient_BfAddMulti(t *testing.T) {
	client.FlushAll()
	ret, err := client.BfAddMulti("test_add_multi", []string{"a", "b", "c"})
	assert.Nil(t, err)
	assert.NotNil(t, ret)
}

func TestClient_BfExistsMulti(t *testing.T) {
	client.FlushAll()
	key := "test_exists_multi"
	ret, err := client.BfAddMulti(key, []string{"a", "b", "c"})
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	existsResult, err := client.BfExistsMulti(key, []string{"a", "b", "notexists"})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(existsResult))
	assert.Equal(t, int64(1), existsResult[0])
	assert.Equal(t, int64(1), existsResult[1])
	assert.Equal(t, int64(0), existsResult[2])
}

func TestClient_TopkReserve(t *testing.T) {
	client.FlushAll()
	ret, err := client.TopkReserve("test_topk_reserve", 10, 2000, 7, 0.925)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
}

func TestClient_TopkAdd(t *testing.T) {
	client.FlushAll()
	key := "test_topk_add"
	ret, err := client.TopkReserve(key, 10, 2000, 7, 0.925)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	rets, err := client.TopkAdd(key, []string{"test", "test1", "test3"})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(rets))
}

func TestClient_TopkQuery(t *testing.T) {
	client.FlushAll()
	key := "test_topk_query"
	ret, err := client.TopkReserve(key, 10, 2000, 7, 0.925)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	rets, err := client.TopkAdd(key, []string{"test"})
	assert.Nil(t, err)
	assert.NotNil(t, rets)
	queryRet, err := client.TopkQuery(key, []string{"test", "nonexist"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(queryRet))
	assert.Equal(t, int64(1), queryRet[0])
	assert.Equal(t, int64(0), queryRet[1])

	key1 := "test_topk_list"
	ret, err = client.TopkReserve(key1, 3, 50, 3, 0.9)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	client.TopkAdd(key1, []string{"A", "B", "C", "D", "E", "A", "A", "B", "C",
		"G", "D", "B", "D", "A", "E", "E"})
	keys, err := client.TopkList(key1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(keys))
	assert.Equal(t, []string{"D", "A", "B"}, keys)
}
