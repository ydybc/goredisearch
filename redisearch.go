package go_redisearch

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-redisearch/goRedis"
	"go-redisearch/redisearch"
)

type RS struct {
	R         *redis.Client
	Ctx       context.Context
	IndexName string
}

func NewClient(addr, pass string, dbNum, poolSize int, indexName string) (RS, error) {
	client, err := redisearch.InitClient(addr, pass, dbNum, poolSize)
	if err != nil {
		return RS{}, err
	}
	return RS{R: client, Ctx: context.Background(), IndexName: indexName}, nil
}

func (r RS) SetIndex(name string) {
	r.IndexName = name
}

func (r RS) DropIndex(deleteDocuments bool) (err error) {
	if deleteDocuments {
		_, err = redisearch.ClientDo(r.R, r.Ctx, "FT.DROPINDEX", r.IndexName, "DD")
	} else {
		_, err = redisearch.ClientDo(r.R, r.Ctx, "FT.DROPINDEX", r.IndexName)
	}
	return
}
func (r RS) CreateIndexWithIndexDefinition(schema *redisearch.Schema, definition *redisearch.IndexDefinition) (err error) {
	return r.indexWithDefinition(schema, definition)
}
func NewDocument(id string, score float32) redisearch.Document {
	return redisearch.NewDocument(id, score)
}
func (r RS) Index(docs ...redisearch.Document) error {
	return redisearch.Index(r.R, r.Ctx, docs...)
}
func (r RS) Search(q *redisearch.Query) (docs []redisearch.Document, total int, err error) {
	return redisearch.Search(r.R, r.Ctx, r.IndexName, q)
}

// internal method
func (r RS) indexWithDefinition(schema *redisearch.Schema, definition *redisearch.IndexDefinition) (err error) {
	args := goRedis.Args{"FT.CREATE", r.IndexName}
	if definition != nil {
		args = definition.Serialize(args)
	}
	// Set flags based on options
	args, err = redisearch.SerializeSchema(schema, args)
	if err != nil {
		return
	}
	//conn := i.pool.Get()
	//defer conn.Close()
	_, err = redisearch.ClientDo(r.R, r.Ctx, args...)

	//_, err = conn.Do("FT.CREATE", args...)
	return
}
