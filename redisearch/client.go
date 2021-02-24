package redisearch

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ydybc/go-redisearch/goRedis"
	"log"
	"time"
)

var (
	//rdb  *redis.Client
	//rCtx context.Context = context.Background()
	NotFound = redis.Nil
)

func InitClient(addr, pass string, dbNum, poolSize int) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       0,    // use default DB
		PoolSize: 1000, // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	return rdb, err
}
func ClientDo(rdb *redis.Client, ctx context.Context, args ...interface{}) (interface{}, error) {
	return rdb.Do(ctx, args...).Result()
}
func Search(rdb *redis.Client, ctx context.Context, name string, q *Query) (docs []Document, total int, err error) {
	args := goRedis.Args{"FT.SEARCH", name}
	args = append(args, q.serialize()...)
	res, err := goRedis.Values(rdb.Do(ctx, args...).Result())
	if err != nil {
		return
	}
	if total, err = goRedis.Int(res[0], nil); err != nil {
		return
	}
	docs = make([]Document, 0, len(res)-1)
	skip := 1
	scoreIdx := -1
	fieldsIdx := -1
	payloadIdx := -1
	if q.Flags&QueryWithScores != 0 {
		scoreIdx = 1
		skip++
	}
	if q.Flags&QueryWithPayloads != 0 {
		payloadIdx = skip
		skip++
	}
	if q.Flags&QueryNoContent == 0 {
		fieldsIdx = skip
		skip++
	}
	if len(res) > skip {
		for i := 1; i < len(res); i += skip {
			if d, e := loadDocument(res, i, scoreIdx, payloadIdx, fieldsIdx); e == nil {
				docs = append(docs, d)
			} else {
				log.Print("Error parsing doc: ", e)
			}
		}
	}
	return
}
