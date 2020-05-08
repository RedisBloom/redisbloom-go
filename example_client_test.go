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
