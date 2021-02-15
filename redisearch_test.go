package go_redisearch

import (
	"go-redisearch/redisearch"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	rs, err := NewClient("192.168.1.117:6666", "", 0, 1000, "testIndex")
	if err!=nil{
		t.Error("NewClient",err)
	}
	//Weight 设置权重
	title:=redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Weight: 10.0, Sortable: true})
	body:=redisearch.NewTextFieldOptions("body", redisearch.TextFieldOptions{Weight: 2.0, Sortable: true})

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(body).
		AddField(title).
		AddField(redisearch.NewNumericField("date"))
		//设置语言
	i:=redisearch.NewIndexDefinition().SetLanguage("chinese").SetLanguageField("chinese")
	rs.DropIndex(true)
	if err := rs.CreateIndexWithIndexDefinition(sc,i); err != nil {
		t.Error("CreateIndex",err)
	}
	// Create a document with an id and given score
	doc1 := NewDocument("doc1",0.3).
		Set("title", "你真好").
		Set("body", "还是这个样子").
		Set("date", time.Now().Unix())
	doc2 := NewDocument("doc2",0.3).
		Set("title", "你不太好").
		Set("body", "还是这个样子").
		Set("date", time.Now().Unix())
	doc3 := NewDocument("doc3",0.3).
		Set("title", "你不太好").
		Set("body", "你真好").
		Set("date", time.Now().Unix())
	doc4 := NewDocument("doc4",0.3).
		Set("title", "还是这个样子").
		Set("body", "你的样子不太好").
		Set("date", time.Now().Unix())
	if err := rs.Index([]redisearch.Document{doc1,doc2,doc3,doc4}...); err != nil {
		t.Error("Index",err)
	}
	// Searching with limit and sorting
	docs, total, err := rs.Search(redisearch.NewQuery("你").
		SetFlags(redisearch.QueryWithScores). //显示评分
		SetLanguage("chinese"). //使用什么分词器
		SetInFields("title","body"). //在什么字段内搜索
		Limit(2, 2))
	t.Log(docs,total,err)
}
