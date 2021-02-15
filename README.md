
# RediSearch Go Client
Go client for RediSearch based on [go-redis](https://github.com/go-redis/redis)

简单易用高性能的 [RediSearch](http://redisearch.io) 全文索引Go库, 基于 [go-redis](https://github.com/go-redis/redis)

轻量简单告别es,

##  Intor
github:https://github.com/RediSearch/RediSearch/

redislabs:http://redisearch.io

Golang [RediSearch](http://redisearch.io) 库是基于 [redigo](https://github.com/gomodule/redigo) 的 [redisearch-go](https://github.com/RediSearch/redisearch-go),

但是更新较慢.好多RediSearch2.0的命令并不支持

### 暂时只实现了基本功能api 
只有 创建 添加&更新 查询 删除 其他等后续如果需要再添加

go-redisearch  | RediSearch  | 介绍
 ---- | ----- | ------  
 CreateIndex  | [FT.CREATE](https://oss.redislabs.com/redisearch/Commands.html#ftcreate) | *
 Index  | [HSET](https://oss.redislabs.com/redisearch/Commands/#hsethsetnxhdelhincrbyhdecrby) | **FT.ADD** 已经被RediSearch2.0淘汰,改为 **HSET**
 Search  | [FT.SEARCH](https://oss.redislabs.com/redisearch/Commands.html#ftsearch) | * 
 DropIndex  | [FT.DROPINDEX](https://oss.redislabs.com/redisearch/Commands/#ftdropindex) | **FT.DROP** 也已经被淘汰,改为 **FT.DROPINDEX**
 
 ### 感谢
 引用借鉴了[redigo](https://github.com/gomodule/redigo) , [redisearch-go](https://github.com/RediSearch/redisearch-go)
