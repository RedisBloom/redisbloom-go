[![license](https://img.shields.io/github/license/RedisBloom/redisbloom-go.svg)](https://github.com/RedisBloom/redisbloom-go)
[![CircleCI](https://circleci.com/gh/RedisBloom/redisbloom-go.svg?style=svg)](https://circleci.com/gh/RedisBloom/redisbloom-go)
[![GitHub issues](https://img.shields.io/github/release/RedisBloom/redisbloom-go.svg)](https://github.com/RedisBloom/redisbloom-go/releases/latest)
[![Codecov](https://codecov.io/gh/RedisBloom/redisbloom-go/branch/master/graph/badge.svg)](https://codecov.io/gh/RedisBloom/redisbloom-go)
[![GoDoc](https://godoc.org/github.com/RedisBloom/redisbloom-go?status.svg)](https://godoc.org/github.com/RedisBloom/redisbloom-go)

# redisbloom-go
[![Forum](https://img.shields.io/badge/Forum-RedisBloom-blue)](https://forum.redislabs.com/c/modules/redisbloom)
[![Discord](https://img.shields.io/discord/697882427875393627?style=flat-square)](https://discord.gg/wXhwjCQ)

Go client for RedisBloom (https://github.com/RedisBloom/redisbloom), based on redigo.

## Installing

```sh
$ go get github.com/RedisBloom/redisbloom-go
```

## Running tests

A simple test suite is provided, and can be run with:

```sh
$ go test
```

The tests expect a Redis server with the RedisBloom module loaded to be available at localhost:6379. You can easily launch RedisBloom with Docker in the following manner:
```
docker run -d -p 6379:6379 --name redis-redisbloom redislabs/rebloom:latest 
```

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
| [BF.MADD](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfmadd) | [BfAddMulti](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.BfAddMulti)  |
| [BF.INSERT](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfinsert) | [BfInsert](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.BfInsert) |
| [BF.EXISTS](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfexists) | [Exists](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Exists) |
| [BF.MEXISTS](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfmexists) | [BfExistsMulti](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.BfExistsMulti) |
| [BF.SCANDUMP](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfscandump) | [BfScanDump](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.BfScanDump) |
| [BF.LOADCHUNK](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfloadchunk) | [BfLoadChunk](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.BfLoadChunk) |
| [BF.INFO](https://oss.redislabs.com/redisbloom/Bloom_Commands/#bfinfo) | [Info](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.Info) |

### Cuckoo Filter

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [CF.RESERVE](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfreserve) | [CfReserve](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfReserve) |
| [CF.ADD](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfadd) |  [CfAdd](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfAdd) |
| [CF.ADDNX](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfaddnx) |  [CfAddNx](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfAddNx) |
| [CF.INSERT](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinsert) |  [CfInsert](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfInsert) |
| [CF.INSERTNX](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinsertnx) |  [CfInsertNx](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfInsertNx) |
| [CF.EXISTS](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfexists) |  [CfExists](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfExists) |
| [CF.DEL](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfdel) |  [CfDel](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfDel) |
| [CF.COUNT](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfcount) |  [CfCount](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfCount) |
| [CF.SCANDUMP](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfscandump) | [CfScanDump](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfScanDump) |
| [CF.LOADCHUNK](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfloadchunck) |  [CfLoadChunk](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfLoadChunk) |
| [CF.INFO](https://oss.redislabs.com/redisbloom/Cuckoo_Commands/#cfinfo) |  [CfInfo](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CfInfo) |

### Count-Min Sketch

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [CMS.INITBYDIM](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinitbydim) | [CmsInitByDim](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsInitByDim) |
| [CMS.INITBYPROB](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinitbyprob) |  [CmsInitByProb](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsInitByProb) |
| [CMS.INCRBY](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsincrby) |  [CmsIncrBy](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsIncrBy) |
| [CMS.QUERY](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsquery) | [CmsQuery](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsQuery) |
| [CMS.MERGE](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsmerge) |  [CmsMerge](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsMerge) |
| [CMS.INFO](https://oss.redislabs.com/redisbloom/CountMinSketch_Commands/#cmsinfo) |  [CmsInfo](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.CmsInfo) |

### TopK Filter

| Command | Recommended API and godoc  |
| :---          |  ----: |
| [TOPK.RESERVE](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkreserve) |  [TopkReserve](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkReserve)  |
| [TOPK.ADD](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkadd) |   [TopkAdd](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkAdd)  |
| [TOPK.INCRBY](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkincrby) |  [TopkIncrby](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkIncrby)  |
| [TOPK.QUERY](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkquery) |   [TopkQuery](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkQuery)  |
| [TOPK.COUNT](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkcount) |   [TopkCount](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkCount)  |
| [TOPK.LIST](https://oss.redislabs.com/redisbloom/TopK_Commands/#topklist) |   [TopkList](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkList)  |
| [TOPK.INFO](https://oss.redislabs.com/redisbloom/TopK_Commands/#topkinfo) |   [TopkInfo](https://godoc.org/github.com/RedisBloom/redisbloom-go#Client.TopkInfo)  |


## License

redisbloom-go is distributed under the BSD 3-Clause license - see [LICENSE](LICENSE)
