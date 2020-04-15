[![license](https://img.shields.io/github/license/RedisBloom/redisbloom-go.svg)](https://github.com/RedisBloom/redisbloom-go)
[![CircleCI](https://circleci.com/gh/RedisBloom/redisbloom-go.svg?style=svg)](https://circleci.com/gh/RedisBloom/redisbloom-go)
[![GitHub issues](https://img.shields.io/github/release/RedisBloom/redisbloom-go.svg)](https://github.com/RedisBloom/redisbloom-go/releases/latest)
[![Codecov](https://codecov.io/gh/RedisBloom/redisbloom-go/branch/master/graph/badge.svg)](https://codecov.io/gh/RedisBloom/redisbloom-go)
[![GoDoc](https://godoc.org/github.com/RedisBloom/redisbloom-go?status.svg)](https://godoc.org/github.com/RedisBloom/redisbloom-go)


# redisbloom-go

Go client for RedisBloom (https://github.com/RedisBloom/redisbloom), based on redigo.

## Installing

```sh
$ go get github.com/RedisBloom/redisbloom-go
```

## Running tests

A simple test suite is provided, and can be run with:

```sh
$ RedisBloom_TEST_PASSWORD="" go test
```

The tests expect a Redis server with the RedisBloom module loaded to be available at localhost:6379

## Example Code

```go
package main 

import (
        "fmt"
        RedisBloom "github.com/RedisBloom/redisbloom-go"
)

func main() {
		// Connect to localhost with no password
        var client = RedisBloom.NewClient("localhost:6379", "nohelp", nil)
        var keyname = "mytest"
        _, haveit := client.Info(keyname)
        if haveit != nil {
			client.CreateKeyWithOptions(keyname, RedisBloom.DefaultCreateOptions)
			client.CreateKeyWithOptions(keyname+"_avg", RedisBloom.DefaultCreateOptions)
			client.CreateRule(keyname, RedisBloom.AvgAggregation, 60, keyname+"_avg")
        }
		// Add sample with timestamp from server time and value 100
        // TS.ADD mytest * 100 
        _, err := client.AddAutoTs(keyname, 100)
        if err != nil {
                fmt.Println("Error:", err)
        }
}
```

## Supported RedisBloom Commands

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [BF.ADD](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfadd) |  |
| [BF.EXISTS](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfexists) |  |

## License

redisbloom-go is distributed under the BSD 3-Clause license - see [LICENSE](LICENSE)
