package redis_bloom_go

import (
	"os"
//	"reflect"
	"testing"
	"time"

//	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func createClient() *Client {
	valueh, exists := os.LookupEnv("REDISBLOOM_TEST_HOST")
	host := "localhost:6379"
	if exists && valueh != "" {
		host = valueh
	}
	valuep, exists := os.LookupEnv("REDISBLOOM_TEST_PASSWORD")
	password := "SUPERSECRET"
	var ptr *string = nil
	if exists {
		password = valuep
	}
	if len(password) > 0 {
		ptr = &password
	}
	return NewClient(host, "test_client", ptr)
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
			"Capacity": 1000,
			"Expansion rate": 2,
            "Number of filters": 1,
            "Number of items inserted": 0,
            "Size": 932,
	})
	
	err = client.Reserve(key, 0.1, 1000)
	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	client.FlushAll()
	key := "test_ADD"
	value := "test_ADD_value";
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
