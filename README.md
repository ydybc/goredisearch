
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
新增 自动补全

搜索 [example](https://github.com/ydybc/go-redisearch/blob/master/redisearch_test.go)  | RediSearch  | other
---- | ----- | ------  
CreateIndex  | [FT.CREATE](https://oss.redislabs.com/redisearch/Commands.html#ftcreate) | *
Set  | [HSET](https://oss.redislabs.com/redisearch/Commands/#hsethsetnxhdelhincrbyhdecrby) |RS2支持**HSET**,可直接使用redis方法操作cudr 
Search  | [FT.SEARCH](https://oss.redislabs.com/redisearch/Commands.html#ftsearch) | * 
DropIndex  | [FT.DROPINDEX](https://oss.redislabs.com/redisearch/Commands/#ftdropindex) | **FT.DROP** 也已经被淘汰,改为 **FT.DROPINDEX**

自动补全 [example](https://github.com/ydybc/go-redisearch/blob/master/auto_conplete_test.go) | RediSearch  | other
---- | ----- | ------  
AddSuggestions  | [FT.SUGADD](https://oss.redislabs.com/redisearch/Commands.html#ftsugadd) | 见test example
GetSuggestions  |[FT.SUGGET](https://oss.redislabs.com/redisearch/Commands.html#ftsugget) |
DeleteSuggestions  | [FT.SUGDEL](https://oss.redislabs.com/redisearch/Commands.html#ftsugdel) |
Length  | [FT.SUGLEN](https://oss.redislabs.com/redisearch/Commands.html#ftsuglen) | 

 
 ### 感谢
 引用借鉴了[redigo](https://github.com/gomodule/redigo) , [redisearch-go](https://github.com/RediSearch/redisearch-go)
