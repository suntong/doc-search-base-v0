package main

import (
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/leopku/bleve-gse-tokenizer"
)

func main() {
	INDEX_DIR := "bleve.gse"
	message := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"

	mapping := bleve.NewIndexMapping()
	os.RemoveAll(INDEX_DIR)
	defer os.RemoveAll(INDEX_DIR)

	if err := mapping.AddCustomTokenizer("gse", map[string]interface{}{
		"type":       "gse",
		"user_dicts": "",
	}); err != nil {
		panic(err)
	}
	if err := mapping.AddCustomAnalyzer("gse", map[string]interface{}{
		"type":      "gse",
		"tokenizer": "gse",
	}); err != nil {
		panic(err)
	}
	mapping.DefaultAnalyzer = "gse"

	index, err := bleve.New(INDEX_DIR, mapping)
	if err != nil {
		panic(err)
	}
	if err := index.Index("1", message); err != nil {
		panic(err)
	}

	query := "干事亲口交待"
	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	req.Highlight = bleve.NewHighlight()
	res, err := index.Search(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result of: '%s': %d matches\n", query, res.Total)
	for i, hit := range res.Hits {
		rv := fmt.Sprintf("%d. %s, (%f)\n", i+res.Request.From+1, hit.ID, hit.Score)
		for fragmentField, fragments := range hit.Fragments {
			rv += fmt.Sprintf("%s: ", fragmentField)
			for _, fragment := range fragments {
				rv += fmt.Sprintf("%s", fragment)
			}
		}
		fmt.Printf("%s\n", rv)
	}

	index.Close()
}
