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
| [TS.CREATE](https://oss.redislabs.com/RedisBloom/commands/#tscreate) |   [CreateKeyWithOptions](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CreateKeyWithOptions)          |
| [TS.ALTER](https://oss.redislabs.com/RedisBloom/commands/#tsalter) |   N/A          |
| [TS.ADD](https://oss.redislabs.com/RedisBloom/commands/#tsadd) |   <ul><li>[Add](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Add)</li><li>[AddAutoTs](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.AddAutoTs)</li><li>[AddWithOptions](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.AddWithOptions)</li><li>[AddAutoTsWithOptions](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.AddWithOptions)</li> </ul>          |
| [TS.MADD](https://oss.redislabs.com/RedisBloom/commands/#tsmadd) |    N/A |
| [TS.INCRBY/TS.DECRBY](https://oss.redislabs.com/RedisBloom/commands/#tsincrbytsdecrby) |    N/A         |
| [TS.CREATERULE](https://oss.redislabs.com/RedisBloom/commands/#tscreaterule) |   [CreateRule](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CreateRule)          |
| [TS.DELETERULE](https://oss.redislabs.com/RedisBloom/commands/#tsdeleterule) |   [DeleteRule](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.DeleteRule)          |
| [TS.RANGE](https://oss.redislabs.com/RedisBloom/commands/#tsrange) |   [RangeWithOptions](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.RangeWithOptions)          |
| [TS.MRANGE](https://oss.redislabs.com/RedisBloom/commands/#tsmrange) |   [MultiRangeWithOptions](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.MultiRangeWithOptions)          |
| [TS.GET](https://oss.redislabs.com/RedisBloom/commands/#tsget) |   [Get](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Get)          |
| [TS.MGET](https://oss.redislabs.com/RedisBloom/commands/#tsmget) |   [MultiGet](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.MultiGet)          |
| [TS.INFO](https://oss.redislabs.com/RedisBloom/commands/#tsinfo) |   [Info](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Info)          |
| [TS.QUERYINDEX](https://oss.redislabs.com/RedisBloom/commands/#tsqueryindex) |    N/A |


## License

redisbloom-go is distributed under the Apache-2 license - see [LICENSE](LICENSE)
