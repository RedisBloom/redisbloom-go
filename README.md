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
        redisbloom "github.com/RedisBloom/redisbloom-go"
)

func main() {
		// Connect to localhost with no password
    var client = redisbloom.NewClient("localhost:6379", "nohelp", nil)
       
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
}
```

## Supported RedisBloom Commands

### Bloom Filter

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [BF.RESERVE](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfreserve) | [Reserve](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Reserve) |
| [BF.ADD](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfadd) | [Add](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Add) |
| [BF.MADD](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfmadd) | N/A |
| [BF.INSERT](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfinsert) | N/A |
| [BF.EXISTS](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfexists) | [Exists](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Exists) |
| [BF.MEXISTS](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfmexists) | N/A |
| [BF.SCANDUMP](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfscandump) | N/A |
| [BF.LOADCHUNK](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfloadchunk) | N/A |
| [BF.INFO](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfinfo) | [Info](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Info) |

### Cuckoo Filter

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [CF.RESERVE](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfreserve) | N/A |
| [CF.ADD](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfadd) |  N/A |
| [CF.ADDNX](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfaddnx) |  N/A |
| [CF.INSERT](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinsert) |  N/A |
| [CF.INSERTNX](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinsertnx) |  N/A |
| [CF.EXISTS](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfexists) |  N/A |
| [CF.DEL](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfdel) |  N/A |
| [CF.COUNT](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfcount) |  N/A |
| [CF.SCANDUMP](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfscandump) |  N/A |
| [CF.LOADCHUNK](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfloadchunck) |  N/A |
| [CF.INFO](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinfo) |  N/A |

### Count-Min Sketch

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [CMS.INITBYDIM](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinitbydim) | N/A |
| [CMS.INITBYPROB](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinitbyprob) |  N/A |
| [CMS.INCRBY](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsincrby) |  N/A |
| [CMS.QUERY](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsquery) |  N/A |
| [CMS.MERGE](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsmerge) |  N/A |
| [CMS.INFO](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinfo) |  N/A |

### TopK Filter

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [TOPK.RESERVE](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkreserve) | N/A |
| [TOPK.ADD](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkadd) |  N/A |
| [TOPK.INCRBY](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkincrby) |  N/A |
| [TOPK.QUERY](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkquery) |  N/A |
| [TOPK.COUNT](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkcount) |  N/A |
| [TOPK.LIST](https://oss.redislabs.com/redisbloom/TopK_Commands/#topklist) |  N/A |
| [TOPK.INFO](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkinfo) |  N/A |


## License

redisbloom-go is distributed under the BSD 3-Clause license - see [LICENSE](LICENSE)
