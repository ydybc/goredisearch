package go_redisearch

import (
	"context"
	"github.com/ydybc/go-redisearch/redisearch"
)

//自动补全
func NewAutoCompleter(addr, pass string, dbNum, poolSize int, indexName string) (redisearch.Autocompleter, error) {
	client, err := redisearch.InitClient(addr, pass, dbNum, poolSize)
	if err != nil {
		return redisearch.Autocompleter{}, err
	}
	return redisearch.Autocompleter{R: client, Ctx: context.Background(), IndexName: indexName}, nil
}
func NewEmptySugList() []redisearch.Suggestion {
	return []redisearch.Suggestion{}
}
