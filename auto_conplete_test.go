package go_redisearch

import (
	"go-redisearch/redisearch"
	"log"
	"testing"
)

//测试自动补全
func TestNewAutoCompleter(t *testing.T) {
	ac, err := NewAutoCompleter("192.168.1.117:6666", "", 0, 1000, "testAutoC")
	if err != nil {
		t.Error("NewAutoCompleter", err)
	}
	//
	//初始化需要补全的句子 OR var sugs []redisearch.Suggestion{}
	sugs := NewEmptySugList()
	sugs = append(sugs, redisearch.Suggestion{Term: "还挺好", Score: 0.5})
	sugs = append(sugs, redisearch.Suggestion{Term: "还不错", Score: 0.5})
	sugs = append(sugs, redisearch.Suggestion{Term: "还挺好还不错", Score: 0.5})
	sugs = append(sugs, redisearch.Suggestion{Term: "还不错还挺好", Score: 0.5})
	sugs = append(sugs, redisearch.Suggestion{Term: "还早", Score: 0.5})
	sugs = append(sugs, redisearch.Suggestion{Term: "早上好", Score: 0.5})
	//插入补全句
	err = ac.AddTerms(sugs...)
	if err != nil {
		log.Fatal("AddTerms", err)
	}
	//查看补全
	opts, err := ac.SuggestOpts("还", redisearch.SuggestOptions{Num: 10, Fuzzy: true, WithScores: true})
	if err != nil {
		log.Fatal("SuggestOpts", err)
	}
	t.Logf("查询补全句%+v\n", opts)
	//删除一些
	err = ac.DeleteTerms(sugs[0:2]...)
	if err != nil {
		log.Fatal("DeleteTerms", err)
	}
	//删除后查看补全
	opts, err = ac.SuggestOpts("还", redisearch.SuggestOptions{Num: 10, Fuzzy: true, WithScores: true})
	if err != nil {
		log.Fatal("SuggestOpts", err)
	}
	t.Logf("删除一些查询%+v\n", opts)
	//删除key
	err = ac.Delete()
	if err != nil {
		log.Fatal("Delete", err)
	}
	/*
		//删除key后查看补全
		opts, err = ac.SuggestOpts("还", redisearch.SuggestOptions{Num: 10, Fuzzy: true,WithScores:true})
		if err!=nil{
			log.Fatal("SuggestOpts",err)
		}
		t.Logf("删除key后查看补全%+v\n",opts)
	*/
}
