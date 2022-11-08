package redis_bloom_go

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// TODO: refactor this hard limit and revise client locking
// Client Max Connections
var maxConns = 500

// Client is an interface to RedisBloom redis commands
type Client struct {
	Pool ConnPool
	Name string
}

// TDigestInfo is a struct that represents T-Digest properties
type TDigestInfo struct {
	compression       int64
	capacity          int64
	mergedNodes       int64
	unmergedNodes     int64
	mergedWeight      int64
	unmergedWeight    int64
	totalCompressions int64
}

// Compression - returns the compression of TDigestInfo instance
func (info *TDigestInfo) Compression() int64 {
	return info.compression
}

// Capacity - returns the capacity of TDigestInfo instance
func (info *TDigestInfo) Capacity() int64 {
	return info.capacity
}

// MergedNodes - returns the merged nodes of TDigestInfo instance
func (info *TDigestInfo) MergedNodes() int64 {
	return info.mergedNodes
}

// UnmergedNodes -  returns the unmerged nodes of TDigestInfo instance
func (info *TDigestInfo) UnmergedNodes() int64 {
	return info.unmergedNodes
}

// MergedWeight - returns the merged weight of TDigestInfo instance
func (info *TDigestInfo) MergedWeight() int64 {
	return info.mergedWeight
}

// UnmergedWeight - returns the unmerged weight of TDigestInfo instance
func (info *TDigestInfo) UnmergedWeight() int64 {
	return info.unmergedWeight
}

// TotalCompressions - returns the total compressions of TDigestInfo instance
func (info *TDigestInfo) TotalCompressions() int64 {
	return info.totalCompressions
}

// NewClient creates a new client connecting to the redis host, and using the given name as key prefix.
// Addr can be a single host:port pair, or a comma separated list of host:port,host:port...
// In the case of multiple hosts we create a multi-pool and select connections at random
// Deprecated: Please use NewClientFromPool() instead
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

// BfAddMulti - Adds one or more items to the Bloom Filter, creating the filter if it does not yet exist.
// args:
// key - the name of the filter
// item - One or more items to add
func (client *Client) BfAddMulti(key string, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err := conn.Do("BF.MADD", args...)
	return redis.Int64s(result, err)
}

// BfExistsMulti - Determines if one or more items may exist in the filter or not.
// args:
// key - the name of the filter
// item - one or more items to check
func (client *Client) BfExistsMulti(key string, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err := conn.Do("BF.MEXISTS", args...)
	return redis.Int64s(result, err)
}

// Begins an incremental save of the bloom filter.
func (client *Client) BfScanDump(key string, iter int64) (int64, []byte, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	reply, err := redis.Values(conn.Do("BF.SCANDUMP", key, iter))
	if err != nil || len(reply) != 2 {
		return 0, nil, err
	}
	iter = reply[0].(int64)
	if reply[1] == nil {
		return iter, nil, err
	}
	return iter, reply[1].([]byte), err
}

// Restores a filter previously saved using SCANDUMP .
func (client *Client) BfLoadChunk(key string, iter int64, data []byte) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("BF.LOADCHUNK", key, iter, data))
}

// This command will add one or more items to the bloom filter, by default creating it if it does not yet exist.
func (client *Client) BfInsert(key string, cap int64, errorRatio float64, expansion int64, noCreate bool, nonScaling bool, items []string) (res []int64, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}
	if cap > 0 {
		args = args.Add("CAPACITY", cap)
	}
	if errorRatio > 0 {
		args = args.Add("ERROR", errorRatio)
	}
	if expansion > 0 {
		args = args.Add("EXPANSION", expansion)
	}
	if noCreate {
		args = args.Add("NOCREATE")
	}
	if nonScaling {
		args = args.Add("NONSCALING")
	}
	args = args.Add("ITEMS").AddFlat(items)
	var resp []interface{}
	var innerRes int64
	resp, err = redis.Values(conn.Do("BF.INSERT", args...))
	if err != nil {
		return
	}
	for _, arrayPos := range resp {
		innerRes, err = redis.Int64(arrayPos, err)
		if err == nil {
			res = append(res, innerRes)
		} else {
			break
		}
	}
	return
}

// Initializes a TopK with specified parameters.
func (client *Client) TopkReserve(key string, topk int64, width int64, depth int64, decay float64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	result, err := conn.Do("TOPK.RESERVE", key, topk, width, depth, strconv.FormatFloat(decay, 'g', 16, 64))
	return redis.String(result, err)
}

// Adds an item to the data structure.
func (client *Client) TopkAdd(key string, items []string) ([]string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err := conn.Do("TOPK.ADD", args...)
	return redis.Strings(result, err)
}

// Returns count for an item.
func (client *Client) TopkCount(key string, items []string) (result []int64, err error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err = redis.Int64s(conn.Do("TOPK.COUNT", args...))
	return
}

// Checks whether an item is one of Top-K items.
func (client *Client) TopkQuery(key string, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err := conn.Do("TOPK.QUERY", args...)
	return redis.Int64s(result, err)
}

// Return full list of items in Top K list.
func (client *Client) TopkListWithCount(key string) (map[string]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return ParseInfoReply(redis.Values(conn.Do("TOPK.LIST", key, "WITHCOUNT")))
}

func (client *Client) TopkList(key string) ([]string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	result, err := conn.Do("TOPK.LIST", key)
	return redis.Strings(result, err)
}

// Returns number of required items (k), width, depth and decay values.
func (client *Client) TopkInfo(key string) (map[string]string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	reply, err := conn.Do("TOPK.INFO", key)
	values, err := redis.Values(reply, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("expects even number of values result")
	}

	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k := values[i].(string)
		switch v := values[i+1].(type) {
		case []byte:
			m[k] = string(values[i+1].([]byte))
			break
		case int64:
			m[k] = strconv.FormatInt(values[i+1].(int64), 10)
		default:
			return nil, fmt.Errorf("unexpected element type for (Ints,String), got type %T", v)
		}
	}
	return m, err
}

// Increase the score of an item in the data structure by increment.
func (client *Client) TopkIncrBy(key string, itemIncrements map[string]int64) ([]string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}
	for k, v := range itemIncrements {
		args = args.Add(k, v)
	}
	reply, err := conn.Do("TOPK.INCRBY", args...)
	return redis.Strings(reply, err)
}

// Initializes a Count-Min Sketch to dimensions specified by user.
func (client *Client) CmsInitByDim(key string, width int64, depth int64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	result, err := conn.Do("CMS.INITBYDIM", key, width, depth)
	return redis.String(result, err)
}

// Initializes a Count-Min Sketch to accommodate requested capacity.
func (client *Client) CmsInitByProb(key string, error float64, probability float64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	result, err := conn.Do("CMS.INITBYPROB", key, error, probability)
	return redis.String(result, err)
}

// Increases the count of item by increment. Multiple items can be increased with one call.
func (client *Client) CmsIncrBy(key string, itemIncrements map[string]int64) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}
	for k, v := range itemIncrements {
		args = args.Add(k, v)
	}
	result, err := conn.Do("CMS.INCRBY", args...)
	return redis.Int64s(result, err)
}

// Returns count for item.
func (client *Client) CmsQuery(key string, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.AddFlat(items)
	result, err := conn.Do("CMS.QUERY", args...)
	return redis.Int64s(result, err)
}

// Merges several sketches into one sketch, stored at dest key
// All sketches must have identical width and depth.
func (client *Client) CmsMerge(dest string, srcs []string, weights []int64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{dest}.Add(len(srcs)).AddFlat(srcs)
	if weights != nil && len(weights) > 0 {
		args = args.Add("WEIGHTS").AddFlat(weights)
	}
	return redis.String(conn.Do("CMS.MERGE", args...))
}

// Returns width, depth and total count of the sketch.
func (client *Client) CmsInfo(key string) (map[string]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return ParseInfoReply(redis.Values(conn.Do("CMS.INFO", key)))
}

// Create an empty cuckoo filter with an initial capacity of {capacity} items.
func (client *Client) CfReserve(key string, capacity int64, bucketSize int64, maxIterations int64, expansion int64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}.Add(capacity)
	if bucketSize > 0 {
		args = args.Add("BUCKETSIZE", bucketSize)
	}
	if maxIterations > 0 {
		args = args.Add("MAXITERATIONS", maxIterations)
	}
	if expansion > 0 {
		args = args.Add("EXPANSION", expansion)
	}
	return redis.String(conn.Do("CF.RESERVE", args...))
}

// Adds an item to the cuckoo filter, creating the filter if it does not exist.
func (client *Client) CfAdd(key string, item string) (bool, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("CF.ADD", key, item))
}

// Adds an item to a cuckoo filter if the item did not exist previously.
func (client *Client) CfAddNx(key string, item string) (bool, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("CF.ADDNX", key, item))
}

// Adds one or more items to a cuckoo filter, allowing the filter to be created with a custom capacity if it does not yet exist.
func (client *Client) CfInsert(key string, cap int64, noCreate bool, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := GetInsertArgs(key, cap, noCreate, items)
	return redis.Int64s(conn.Do("CF.INSERT", args...))
}

// Adds one or more items to a cuckoo filter, allowing the filter to be created with a custom capacity if it does not yet exist.
func (client *Client) CfInsertNx(key string, cap int64, noCreate bool, items []string) ([]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := GetInsertArgs(key, cap, noCreate, items)
	return redis.Int64s(conn.Do("CF.INSERTNX", args...))
}

func GetInsertArgs(key string, cap int64, noCreate bool, items []string) redis.Args {
	args := redis.Args{key}
	if cap > 0 {
		args = args.Add("CAPACITY", cap)
	}
	if noCreate {
		args = args.Add("NOCREATE")
	}
	args = args.Add("ITEMS").AddFlat(items)
	return args
}

// Check if an item exists in a Cuckoo Filter
func (client *Client) CfExists(key string, item string) (bool, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("CF.EXISTS", key, item))
}

// Deletes an item once from the filter.
func (client *Client) CfDel(key string, item string) (bool, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("CF.DEL", key, item))
}

// Returns the number of times an item may be in the filter.
func (client *Client) CfCount(key string, item string) (int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("CF.COUNT", key, item))
}

// Begins an incremental save of the cuckoo filter.
func (client *Client) CfScanDump(key string, iter int64) (int64, []byte, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	reply, err := redis.Values(conn.Do("CF.SCANDUMP", key, iter))
	if err != nil || len(reply) != 2 {
		return 0, nil, err
	}
	iter = reply[0].(int64)
	if reply[1] == nil {
		return iter, nil, err
	}
	return iter, reply[1].([]byte), err
}

// Restores a filter previously saved using SCANDUMP
func (client *Client) CfLoadChunk(key string, iter int64, data []byte) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("CF.LOADCHUNK", key, iter, data))
}

// Return information about key
func (client *Client) CfInfo(key string) (map[string]int64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return ParseInfoReply(redis.Values(conn.Do("CF.INFO", key)))
}

// TdCreate - Allocate the memory and initialize the t-digest
func (client *Client) TdCreate(key string, compression int64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("TDIGEST.CREATE", key, "COMPRESSION", compression))
}

// TdReset - Reset the sketch to zero - empty out the sketch and re-initialize it
func (client *Client) TdReset(key string) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("TDIGEST.RESET", key))
}

// TdAdd - Adds one or more samples to a sketch
func (client *Client) TdAdd(key string, samples map[float64]float64) (string, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	args := redis.Args{key}
	for k, v := range samples {
		args = args.Add(k, v)
	}
	reply, err := conn.Do("TDIGEST.ADD", args...)
	return redis.String(reply, err)
}

// tdMerge - The internal representation of TdMerge. All underlying functions call this one,
// returning its results. It allows us to maintain interfaces.
// see https://redis.io/commands/tdigest.merge/
//
// The default values for compression is 100
func (client *Client) tdMerge(toKey string, compression int64, override bool, numKeys int64, fromKey ...string) (string, error) {
	if numKeys < 1 {
		return "", errors.New("a minimum of one key must be merged")
	}

	conn := client.Pool.Get()
	defer conn.Close()
	overidable := ""
	if override {
		overidable = "1"
	}
	return redis.String(conn.Do("TDIGEST.MERGE", toKey,
		strconv.FormatInt(numKeys, 10),
		strings.Join(fromKey, " "),
		"COMPRESSION", compression,
		overidable))
}

// TdMerge - Merges all of the values from 'from' to 'this' sketch
func (client *Client) TdMerge(toKey string, numKeys int64, fromKey ...string) (string, error) {
	return client.tdMerge(toKey, 100, false, numKeys, fromKey...)
}

// TdMergeWithCompression - Merges all of the values from 'from' to 'this' sketch with specified compression
func (client *Client) TdMergeWithCompression(toKey string, compression int64, numKeys int64, fromKey ...string) (string, error) {
	return client.tdMerge(toKey, compression, false, numKeys, fromKey...)
}

// TdMergeWithOverride - Merges all of the values from 'from' to 'this' sketch overriding the destination key if it exists
func (client *Client) TdMergeWithOverride(toKey string, override bool, numKeys int64, fromKey ...string) (string, error) {
	return client.tdMerge(toKey, 100, true, numKeys, fromKey...)
}

// TdMergeWithCompressionAndOverride - Merges all of the values from 'from' to 'this' sketch with specified compression
// and overriding the destination key if it exists
func (client *Client) TdMergeWithCompressionAndOverride(toKey string, compression int64, numKeys int64, fromKey ...string) (string, error) {
	return client.tdMerge(toKey, compression, true, numKeys, fromKey...)
}

// TdMin - Get minimum value from the sketch. Will return DBL_MAX if the sketch is empty
func (client *Client) TdMin(key string) (float64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Float64(conn.Do("TDIGEST.MIN", key))
}

// TdMax - Get maximum value from the sketch. Will return DBL_MIN if the sketch is empty
func (client *Client) TdMax(key string) (float64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Float64(conn.Do("TDIGEST.MAX", key))
}

// TdQuantile - Returns an estimate of the cutoff such that a specified fraction of the data added
// to this TDigest would be less than or equal to the cutoff
func (client *Client) TdQuantile(key string, quantile float64) ([]float64, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return redis.Float64s(conn.Do("TDIGEST.QUANTILE", key, quantile))
}

// TdCdf - Returns the list of fractions of all points added which are <= values
func (client *Client) TdCdf(key string, values ...float64) ([]float64, error) {
	conn := client.Pool.Get()
	defer conn.Close()

	args := make([]string, len(values))
	for idx, obj := range values {
		args[idx] = strconv.FormatFloat(obj, 'f', -1, 64)
	}
	return redis.Float64s(conn.Do("TDIGEST.CDF", key, strings.Join(args, " ")))
}

// TdInfo - Returns compression, capacity, total merged and unmerged nodes, the total
// compressions made up to date on that key, and merged and unmerged weight.
func (client *Client) TdInfo(key string) (TDigestInfo, error) {
	conn := client.Pool.Get()
	defer conn.Close()
	return ParseTDigestInfo(redis.Values(conn.Do("TDIGEST.INFO", key)))
}

func ParseInfoReply(values []interface{}, err error) (map[string]int64, error) {
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("expects even number of values result")
	}
	m := make(map[string]int64, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1].(int64)
	}
	return m, err
}

func ParseTDigestInfo(result interface{}, err error) (info TDigestInfo, outErr error) {
	values, outErr := redis.Values(result, err)
	if outErr != nil {
		return TDigestInfo{}, err
	}
	if len(values)%2 != 0 {
		return TDigestInfo{}, errors.New("ParseInfo expects even number of values result")
	}
	var key string
	for i := 0; i < len(values); i += 2 {
		key, outErr = redis.String(values[i], nil)
		if outErr != nil {
			return TDigestInfo{}, outErr
		}
		switch key {
		case "Compression":
			info.compression, outErr = redis.Int64(values[i+1], nil)
		case "Capacity":
			info.capacity, outErr = redis.Int64(values[i+1], nil)
		case "Merged nodes":
			info.mergedNodes, outErr = redis.Int64(values[i+1], nil)
		case "Unmerged nodes":
			info.unmergedNodes, outErr = redis.Int64(values[i+1], nil)
		case "Merged weight":
			info.mergedWeight, outErr = redis.Int64(values[i+1], nil)
		case "Unmerged weight":
			info.unmergedWeight, outErr = redis.Int64(values[i+1], nil)
		case "Total compressions":
			info.totalCompressions, outErr = redis.Int64(values[i+1], nil)
		}
		if outErr != nil {
			return TDigestInfo{}, outErr
		}
	}

	return info, nil
}
