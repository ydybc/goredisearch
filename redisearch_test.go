package goredisearch

import (
	"github.com/ydybc/goredisearch/redisearch"
	"strconv"
	"testing"
	"time"
)

type Dog struct {
	Name    string
	Feature string
	Gender  string
	Length  int //m
	Date    int64
}

//测试中文查询
func TestNewClient(t *testing.T) {
	var (
		instData []redisearch.Document
		dogs     []Dog = []Dog{
			{Name: "超人汪", Feature: "无坚不摧力大无穷的公汪", Gender: "male", Length: 10},
			{Name: "钢铁公汪", Feature: "身上覆盖高科技盔甲的汪汪", Gender: "male", Length: 20},
			{Name: "机械汪", Feature: "科技结晶汪汪不知是公汪还是母汪", Gender: "it", Length: 20},
			{Name: "神奇母汪", Feature: "无坚不摧力大无穷的汪汪", Gender: "female", Length: 15},
		}
	)
	//初始化 link redis
	rs, err := NewSearchClient("192.168.1.117:6666", "", 0, 1000, "testIndex")
	if err != nil {
		t.Error("Client", err)
	}
	//Weight 设置权重
	name := redisearch.NewTextFieldOptions("Name", redisearch.TextFieldOptions{Weight: 2.0, Sortable: true})
	feature := redisearch.NewTextFieldOptions("Feature", redisearch.TextFieldOptions{Weight: 2.0, Sortable: true})
	gender := redisearch.NewTextField("Gender")
	weight := redisearch.NewNumericField("Length")
	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(name).
		AddField(feature).
		AddField(gender).
		AddField(weight).
		AddField(redisearch.NewNumericField("Date")) //时间
	//设置语言
	i := redisearch.NewIndexDefinition().SetLanguage("chinese").SetLanguageField("chinese")
	//删除testIndex索引(如果有的话)
	err = rs.DropIndex(true)
	if err != nil {
		t.Error("DropIndex", err)
	}
	//创建Index
	if err := rs.CreateIndexWithIndexDefinition(sc, i); err != nil {
		t.Error("CreateIndex", err)
	}
	//组合索引数据
	for k, v := range dogs {
		instData = append(instData,
			NewDocument("dogInfo"+strconv.Itoa(k), 1).
				Set("Name", v.Name).
				Set("Feature", v.Feature).
				Set("Gender", v.Gender).
				Set("Length", v.Length).
				Set("Date", time.Now().Unix()))
	}
	//插入&更新 索引 set indexData
	if err = rs.Set(instData...); err != nil {
		t.Error("SetIndex", err)
	}
	// 进行搜索 start Search
	keyWord := "汪"
	docs, total, err := rs.Search(redisearch.NewQuery(keyWord).
		SetFlags(redisearch.QueryWithScores). //评分
		SetLanguage("chinese").               //使用什么分词器
		SetInFields("Name", "Feature").       //在什么字段内搜索
		Limit(0, 4).
		Highlight(nil, "<b>", "</b>"))
	t.Logf("word:%s,res:%+v,total:%d,err:%v\n", keyWord, docs, total, err)
	keyWord = "公汪"
	docs, total, err = rs.Search(redisearch.NewQuery(keyWord).
		SetFlags(redisearch.QueryWithScores).
		SetLanguage("chinese").
		SetInFields("Name", "Feature").
		Limit(0, 4).
		Highlight(nil, "<b>", "</b>"))
	t.Logf("word:%s,res:%+v,total:%d,err:%v\n", keyWord, docs, total, err)
	keyWord = "母汪"
	docs, total, err = rs.Search(redisearch.NewQuery(keyWord).
		SetFlags(redisearch.QueryWithScores).
		SetLanguage("chinese").
		SetInFields("Name", "Feature").
		Limit(0, 4).
		Highlight(nil, "<b>", "</b>"))
	t.Logf("word:%s,res:%+v,total:%d,err:%v\n", keyWord, docs, total, err)
}
