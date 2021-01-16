package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/token/unicodenorm"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/registry"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	var typ, dir, idxPath, query string
	flag.StringVar(&typ, "t", "txt,md", "target file extensions")
	flag.StringVar(&dir, "d", "", "target directory")
	flag.StringVar(&idxPath, "i", "~/.lsh/lsh_index", "lsh index path")
	flag.StringVar(&query, "q", "", "search query")
	flag.Parse()

	var err error
	if idxPath != "" {
		idxPath, err = expandHome(idxPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("-i needs any path")
	}

	println("] Open")
	idx, err := bleve.Open(idxPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		idx, err = createIndex(idxPath)
	}
	defer idx.Close()

	if dir != "" && typ != "" {

		dir, err = expandHome(dir)
		if err != nil {
			log.Fatal(err)
		}

		println("] indexingFiles")
		var types []string
		for _, t := range strings.Split(typ, ",") {
			types = append(types, strings.TrimSpace(t))
		}
		err = indexingFiles(idx, dir, types)
		if err != nil {
			log.Fatal(err)
		}
	}
	if query != "" {
		println("] Query")
		q := bleve.NewQueryStringQuery(query)
		s := bleve.NewSearchRequestOptions(q, 10, 0, true)
		s.Fields = []string{"Path"}
		r, err := idx.Search(s)
		if err != nil {
			log.Fatal(err)
		}
		for _, h := range r.Hits {
			fmt.Printf("%#v\n", h)
		}
	}
}

func expandHome(d string) (string, error) {
	if strings.HasPrefix(d, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, d[1:]), nil
	}
	return d, nil
}

type Doc struct {
	ID      string
	Path    string
	Name    string
	Content string
}

func (d Doc) Type() string {
	return "Doc"
}

var _ mapping.Classifier = (*Doc)(nil)

func indexingFiles(idx bleve.Index, dir string, types []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !isTarget(info.Name(), types) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		bkey := md5.Sum([]byte(path))
		key := hex.EncodeToString(bkey[:])
		d := Doc{
			ID:      key,
			Path:    path,
			Name:    info.Name(),
			Content: string(buf),
		}
		return idx.Index(key, d)
	})
}

func isTarget(file string, exts []string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func createIndex(idxPath string) (bleve.Index, error) {
	m, err := buildMapping()
	if err != nil {
		log.Fatal(err)
	}
	return bleve.New(idxPath, m)
}

func buildMapping() (mapping.IndexMapping, error) {
	// open a new index
	m := bleve.NewIndexMapping()
	dm := bleve.NewDocumentMapping()
	tf := bleve.NewTextFieldMapping()
	dm.AddFieldMappingsAt("Name", tf)
	dm.AddFieldMappingsAt("Content", tf)
	m.AddDocumentMapping("Doc", dm)
	return m, nil
}

const AnalyzerName = "custom"

func AnalyzerConstructor(_ map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}
	normalizeFilter := unicodenorm.MustNewUnicodeNormalizeFilter(unicodenorm.NFKD)
	lowercaseFilter := lowercase.NewLowerCaseFilter()
	ngramFilter := ngram.NewNgramFilter(2, 2)
	rv := analysis.Analyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			normalizeFilter,
			lowercaseFilter,
			ngramFilter,
		},
	}
	return &rv, nil
}

func init() {
	registry.RegisterAnalyzer(AnalyzerName, AnalyzerConstructor)
}
