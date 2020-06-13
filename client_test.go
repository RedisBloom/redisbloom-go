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

func TestClient_CmsInitByDim(t *testing.T) {
	client.FlushAll()
	ret, err := client.CmsInitByDim("test_cms_initbydim", 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
}

func TestClient_CmsInitByProb(t *testing.T) {
	client.FlushAll()
	ret, err := client.CmsInitByProb("test_cms_initbyprob", 0.01, 0.01)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
}

func TestClient_CmsIncrBy(t *testing.T) {
	client.FlushAll()
	key := "test_cms_incrby"
	ret, err := client.CmsInitByDim(key, 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	results, err := client.CmsIncrBy(key, map[string]int64{"foo": 5})
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, int64(5), results[0])
}

func TestClient_CmsQuery(t *testing.T) {
	client.FlushAll()
	key := "test_cms_query"
	ret, err := client.CmsInitByDim(key, 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	results, err := client.CmsQuery(key, []string{"notexist"})
	assert.Nil(t, err)
	assert.NotNil(t, 0, results[0])
	_, err = client.CmsIncrBy(key, map[string]int64{"foo": 5})
	assert.Nil(t, err)
	results, err = client.CmsQuery(key, []string{"foo"})
	assert.Nil(t, err)
	assert.Equal(t, int64(5), results[0])
}

func TestClient_CmsMerge(t *testing.T) {
	client.FlushAll()
	ret, err := client.CmsInitByDim("A", 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	ret, err = client.CmsInitByDim("B", 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	ret, err = client.CmsInitByDim("C", 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	client.CmsIncrBy("A", map[string]int64{"foo": 5, "bar": 3, "baz": 9})
	client.CmsIncrBy("B", map[string]int64{"foo": 2, "bar": 3, "baz": 1})
	client.CmsMerge("C", []string{"A", "B"}, nil)
	results, err := client.CmsQuery("C", []string{"foo", "bar", "baz"})
	assert.Equal(t, []int64{7, 6, 10}, results)
}

func TestClient_CmsInfo(t *testing.T) {
	client.FlushAll()
	key := "test_cms_info"
	ret, err := client.CmsInitByDim(key, 1000, 5)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	info, err := client.CmsInfo(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(1000), info["width"])
	assert.Equal(t, int64(5), info["depth"])
	assert.Equal(t, int64(0), info["count"])
}
