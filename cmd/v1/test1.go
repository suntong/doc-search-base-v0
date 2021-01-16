package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
)

func main() {

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
	}

	err = idx.Index("a", map[string]interface{}{
		"desc": "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
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
	fmt.Println("Default fragment size:")
	fmt.Println(sres)

	_, err = bleve.Config.Cache.DefineFragmenter("short", map[string]interface{}{
		"type": "simple",
		"size": 20.0,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = bleve.Config.Cache.DefineHighlighter("custom", map[string]interface{}{
		"type":       "simple",
		"fragmenter": "short",
		"formatter":  "ansi",
	})
	if err != nil {
		log.Fatal(err)
	}

	// search again with custom highlighter using shorter fragments
	sreq2 := bleve.NewSearchRequest(q)
	sreq2.Highlight = bleve.NewHighlightWithStyle("custom")
	sres2, err := idx.Search(sreq2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Shorter fragment size:")
	fmt.Println(sres2)
}