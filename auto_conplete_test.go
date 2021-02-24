package go_redisearch

import (
	"github.com/ydybc/go-redisearch/redisearch"
	"log"
	"testing"
)

//测试自动补全
func TestNewAutoCompleter(t *testing.T) {
	ac, err := NewAutoCompleterClient("192.168.1.117:6666", "", 0, 1000, "testAutoC")
	if err != nil {
		t.Error("NewAutoCompleterClient", err)
	}
	//
	//初始化需要补全的句子
	sugs := NewEmptySugList() //OR var sugs []redisearch.Suggestion{}
	sugs = append(sugs, redisearch.Suggestion{Term: "还挺好还不错", Score: 1, Incr: true})
	sugs = append(sugs, redisearch.Suggestion{Term: "还不错还挺好", Score: 1, Incr: true})
	sugs = append(sugs, redisearch.Suggestion{Term: "还早", Score: 1, Incr: false})
	sugs = append(sugs, redisearch.Suggestion{Term: "还好", Score: 1, Incr: false})
	sugs = append(sugs, redisearch.Suggestion{Term: "还挺好", Score: 1, Incr: true})
	sugs = append(sugs, redisearch.Suggestion{Term: "还不错", Score: 1, Incr: true})
	sugs = append(sugs, redisearch.Suggestion{Term: "还差不多", Score: 0.5, Incr: true})
	sugs = append(sugs, redisearch.Suggestion{Term: "早上好", Score: 0.5, Incr: true})
	//插入补全句
	err = ac.AddSuggestions(sugs...)
	if err != nil {
		log.Fatal("AddSuggestions", err)
	}
	//查看补全
	opts, err := ac.GetSuggestions("还", redisearch.SuggestOptions{Num: 5, Fuzzy: true, WithScores: true})
	if err != nil {
		log.Fatal("GetSuggestions", err)
	}
	length, err := ac.Length()
	if err != nil {
		log.Fatal("Length", err)
	}
	t.Logf("len:%d,查询补全句%+v\n", length, opts)
	//删除一些
	err = ac.DeleteSuggestions(sugs[0:2]...)
	if err != nil {
		log.Fatal("DeleteSuggestions", err)
	}
	//删除后查看补全
	opts, err = ac.GetSuggestions("还", redisearch.SuggestOptions{Num: 10, Fuzzy: true, WithScores: true})
	if err != nil {
		log.Fatal("GetSuggestions", err)
	}
	length, err = ac.Length()
	if err != nil {
		log.Fatal("Length", err)
	}
	t.Logf("len:%d,删除一些查询%+v\n", length, opts)
	//删除key
	err = ac.Delete()
	if err != nil {
		log.Fatal("Delete", err)
	}
	length, err = ac.Length()
	if err != nil {
		log.Fatal("Length", err)
	}
	t.Logf("Delete after len:%+v\n", length)
	/*
		//删除key后查看补全
		opts, err = ac.GetSuggestions("还", redisearch.SuggestOptions{Num: 10, Fuzzy: true,WithScores:true})
		if err!=nil{
			log.Fatal("GetSuggestions",err)
		}
		t.Logf("删除key后查看补全%+v\n",opts)
	*/
}
