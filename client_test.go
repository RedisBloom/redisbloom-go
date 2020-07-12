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

func TestClient_BfInsert(t *testing.T) {
	client.FlushAll()
	key := "test_bf_insert"
	key_nocreate := "test_bf_insert_nocreate"
	key_noscaling := "test_bf_insert_noscaling"

	ret, err := client.BfInsert(key, 1000, 0.1, -1, false, false, []string{"a"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	assert.True(t, ret[0] > 0)
	existsResult, err := client.BfExistsMulti(key, []string{"a"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(existsResult))
	assert.Equal(t, int64(1), existsResult[0])

	ret, err = client.BfInsert(key, 1000, 0.1, -1, false, false, []string{"a", "b"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ret))

	// Test for NOCREATE : If specified, indicates that the filter should not be created if it does not already exist
	_, err = client.BfInsert(key_nocreate, 1000, 0.1, -1, true, false, []string{"a"})
	assert.NotNil(t, err)

	// Test NONSCALING : Prevents the filter from creating additional sub-filters if initial capacity is reached.
	ret, err = client.BfInsert(key_noscaling, 2, 0.1, -1, false, true, []string{"a", "b"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ret))
	ret, err = client.BfInsert(key_noscaling, 2, 0.1, -1, false, true, []string{"c"})
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(ret))
	assert.Equal(t, err.Error(), "ERR non scaling filter is full")
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

func TestClient_TopkInfo(t *testing.T) {
	client.FlushAll()
	key := "test_topk_info"
	ret, err := client.TopkReserve(key, 10, 2000, 7, 0.925)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)

	info, err := client.TopkInfo(key)
	assert.Equal(t, "10", info["k"])
	assert.Equal(t, "2000", info["width"])
	assert.Equal(t, "7", info["depth"])

	info, err = client.TopkInfo("notexists")
	assert.NotNil(t, err)
}

func TestClient_TopkIncrBy(t *testing.T) {
	client.FlushAll()
	key := "test_topk_incrby"
	ret, err := client.TopkReserve(key, 50, 2000, 7, 0.925)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)

	rets, err := client.TopkAdd(key, []string{"foo", "bar", "42"})
	assert.Nil(t, err)
	assert.NotNil(t, rets)

	rets, err = client.TopkIncrBy(key, map[string]int64{"foo": 3, "bar": 2, "42": 30})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(rets))
	assert.Equal(t, "", rets[2])
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

func TestClient_CfReserve(t *testing.T) {
	client.FlushAll()
	key := "test_cf_reserve"
	ret, err := client.CfReserve(key, 1000, -1, -1, -1)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
}

func TestClient_CfAdd(t *testing.T) {
	client.FlushAll()
	key := "test_cf_add"
	ret, err := client.CfAdd(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	ret, err = client.CfAddNx(key, "b")
	assert.Nil(t, err)
	assert.True(t, ret)
}

func TestClient_CfInsert(t *testing.T) {
	client.FlushAll()
	key := "test_cf_insert"
	ret, err := client.CfInsert(key, 1000, false, []string{"a"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	assert.True(t, ret[0] > 0)
	ret, err = client.CfInsertNx(key, 1000, true, []string{"b"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	assert.True(t, ret[0] > 0)
}

func TestClient_CfExists(t *testing.T) {
	client.FlushAll()
	key := "test_cf_exists"
	ret, err := client.CfAdd(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	ret, err = client.CfExists(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
}

func TestClient_CfDel(t *testing.T) {
	client.FlushAll()
	key := "test_cf_del"
	ret, err := client.CfAdd(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	ret, err = client.CfExists(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	ret, err = client.CfDel(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	ret, err = client.CfExists(key, "a")
	assert.Nil(t, err)
	assert.False(t, ret)
}

func TestClient_CfCount(t *testing.T) {
	client.FlushAll()
	key := "test_cf_count"
	ret, err := client.CfAdd(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	count, err := client.CfCount(key, "a")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestClient_CfScanDump(t *testing.T) {
	client.FlushAll()
	key := "test_cf_scandump"
	ret, err := client.CfReserve(key, 100, 50, -1, -1)
	assert.Nil(t, err)
	assert.Equal(t, "OK", ret)
	client.CfAdd(key, "a")
	curIter := int64(0)
	chunks := make([]map[string]interface{}, 0)
	for {
		iter, data, err := client.CfScanDump(key, curIter)
		assert.Nil(t, err)
		curIter = iter
		if iter == int64(0) {
			break
		}
		chunk := map[string]interface{}{"iter": iter, "data": data}
		chunks = append(chunks, chunk)
	}
	client.FlushAll()
	for i := 0; i < len(chunks); i++ {
		ret, err := client.CfLoadChunk(key, chunks[i]["iter"].(int64), chunks[i]["data"].([]byte))
		assert.Nil(t, err)
		assert.Equal(t, "OK", ret)
	}
	exists, err := client.CfExists(key, "a")
	assert.True(t, exists)
}

func TestClient_CfInfo(t *testing.T) {
	client.FlushAll()
	key := "test_cf_info"
	ret, err := client.CfAdd(key, "a")
	assert.Nil(t, err)
	assert.True(t, ret)
	info, err := client.CfInfo(key)
	assert.Equal(t, int64(1080), info["Size"])
	assert.Equal(t, int64(512), info["Number of buckets"])
	assert.Equal(t, int64(0), info["Number of filter"])
	assert.Equal(t, int64(1), info["Number of items inserted"])
	assert.Equal(t, int64(0), info["Max iteration"])
}

func TestClient_BfScanDump(t *testing.T) {
	client.FlushAll()
	key := "test_bf_scandump"
	err := client.Reserve(key, 0.01, 1000)
	assert.Nil(t, err)
	client.Add(key, "1")
	curIter := int64(0)
	chunks := make([]map[string]interface{}, 0)
	for {
		iter, data, err := client.BfScanDump(key, curIter)
		assert.Nil(t, err)
		curIter = iter
		if iter == int64(0) {
			break
		}
		chunk := map[string]interface{}{"iter": iter, "data": data}
		chunks = append(chunks, chunk)
	}
	client.FlushAll()
	for i := 0; i < len(chunks); i++ {
		ret, err := client.BfLoadChunk(key, chunks[i]["iter"].(int64), chunks[i]["data"].([]byte))
		assert.Nil(t, err)
		assert.Equal(t, "OK", ret)
	}
	exists, err := client.Exists(key, "1")
	assert.True(t, exists)
}
