package redis_bloom_go_test

import (
	"fmt"
	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/gomodule/redigo/redis"
	"log"
)

// exemplifies the NewClient function
func ExampleNewClient() {
	host := "localhost:6379"
	var client = redisbloom.NewClient(host, "nohelp", nil)

	// BF.ADD mytest item
	_, err := client.Add("mytest", "myItem")
	if err != nil {
		fmt.Println("Error:", err)
	}

	exists, err := client.Exists("mytest", "myItem")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("myItem exists in mytest: ", exists)
	// Output: myItem exists in mytest:  true

}

// exemplifies the NewClientFromPool function
func ExampleNewClientFromPool() {
	host := "localhost:6379"
	password := ""
	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", host, redis.DialPassword(password))
	}}
	client := redisbloom.NewClientFromPool(pool, "bloom-client-1")

	// BF.ADD mytest item
	_, err := client.Add("mytest", "myItem")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	exists, err := client.Exists("mytest", "myItem")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("myItem exists in mytest: ", exists)
	// Output: myItem exists in mytest:  true

}

// exemplifies the TdCreate function
func ExampleTdCreate() {
	host := "localhost:6379"
	var client = redisbloom.NewClient(host, "nohelp", nil)

	ret, err := client.TdCreate("key", 100)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(ret)
	// Output: OK

}

// exemplifies the TdCreate function
func ExampleTdAdd() {
	host := "localhost:6379"
	var client = redisbloom.NewClient(host, "nohelp", nil)

	key := "example"
	ret, err := client.TdCreate(key, 100)
	if err != nil {
		fmt.Println("Error:", err)
	}

	samples := map[float64]float64{1.0: 1.0, 2.0: 2.0}
	ret, err = client.TdAdd(key, samples)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(ret)
	// Output: OK

}

// exemplifies the TdMin TdMax and TdCdf functions
func ExampleTdQuery() {
	host := "localhost:6379"
	var client = redisbloom.NewClient(host, "nohelp", nil)

	key := "example"
	_, err := client.TdCreate(key, 10)
	if err != nil {
		fmt.Println("Error:", err)
	}

	samples := map[float64]float64{1.0: 1.0, 2.0: 2.0, 3.0: 3.0}
	_, err = client.TdAdd(key, samples)
	if err != nil {
		fmt.Println("Error:", err)
	}

	min, err := client.TdMin(key)
	if err != nil {
		fmt.Println("Error:", err)
	}

	max, err := client.TdMax(key)
	if err != nil {
		fmt.Println("Error:", err)
	}

	cdf, err := client.TdCdf(key, 0.0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(min, max, cdf)
	// Output: 1 3 0
}
