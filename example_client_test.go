package redis_bloom_go_test

import (
	"fmt"
	redisbloom "github.com/RedisBloom/redisbloom-go"
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
	// myItem exists in mytest:  true

}
