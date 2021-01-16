package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"
)

func main() {

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
	}

	err = idx.Index("a", map[string]interface{}{
		"desc": "Lorem Ipsum is simply dummy text of the printing and typesetting industry. \nLorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. \nIt has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.\nLorem Ipsum is simply dummy text of the printing and typesetting industry. \nLorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer idx.Close()

	q := bleve.NewTermQuery("text")
	sreq := bleve.NewSearchRequest(q)
	sreq.Highlight = bleve.NewHighlightWithStyle("ansi")
	sres, err := idx.Search(sreq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sres)
}
