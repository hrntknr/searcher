package main

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestServiceRegist(t *testing.T) {
	sentenceSplitter := newSentenceSplitterMock([]string{"これはペンです。", "これはりんごです。"})
	charFilterMock := newCharFilterMock([]string{"これはペンです。", "これはりんごです。"})
	tokenizerMock := newTokenizerMock([][]string{{"これ", "ペン", "です"}, {"これ", "りんご", "です"}})
	wordFilterMock := newWordFilterMock([][]string{{"これ", "ペン"}, {"これ", "りんご"}})
	dbMock := newDbMock([]searchResult{})
	service, _ := newService(
		sentenceSplitter,
		tokenizerMock,
		[]charFilter{charFilterMock},
		[]wordFilter{wordFilterMock},
		dbMock,
	)

	err := service.regist("uri", "これはペンです。これはりんごです。")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		[]string{"これはペンです。これはりんごです。"},
		sentenceSplitter.splitArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[][]string{{"これはペンです。", "これはりんごです。"}},
		charFilterMock.fitlerArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[][]string{{"これはペンです。", "これはりんごです。"}},
		tokenizerMock.analyzeArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[][][]string{{{"これ", "ペン", "です"}, {"これ", "りんご", "です"}}},
		wordFilterMock.fitlerArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[]dbStoreArgs{
			{Token: "これ", Positions: []position{{0, 0, "これはペンです。"}, {0, 2, "これはりんごです。"}}, Uri: "uri"},
			{Token: "ペン", Positions: []position{{1, 1, "これはペンです。"}}, Uri: "uri"},
			{Token: "りんご", Positions: []position{{1, 3, "これはりんごです。"}}, Uri: "uri"},
		},
		dbMock.storeArgs,
		cmp.Transformer("Sort", func(in []dbStoreArgs) []dbStoreArgs {
			out := append([]dbStoreArgs(nil), in...)
			sort.Slice(out, func(i int, j int) bool {
				return out[i].Token < out[j].Token
			})
			return out
		}),
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestServiceSearch(t *testing.T) {
	sentenceSplitterMock := newSentenceSplitterMock([]string{})
	tokenizerMock := newTokenizerMock([][]string{{"word1", "word2"}})
	charFilterMock := newCharFilterMock([]string{"body2"})
	wordFilterMock := newWordFilterMock([][]string{{"word3", "word4", "word3"}})
	dbMock := newDbMock([]searchResult{{
		Uri:       "uri",
		Score:     10,
		Sentences: []string{"すもももももももものうち"},
	}})
	service, _ := newService(
		sentenceSplitterMock,
		tokenizerMock,
		[]charFilter{charFilterMock},
		[]wordFilter{wordFilterMock},
		dbMock,
	)

	result, err := service.search("body1", 10, 100)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		[][]string{{"body1"}},
		charFilterMock.fitlerArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[][]string{{"body2"}},
		tokenizerMock.analyzeArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[]dbSearchArgs{{
			Tokens: []string{"word1", "word2"},
			Offset: 10,
			Count:  100,
		}},
		dbMock.searchArgs,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[]searchResult{{Uri: "uri", Score: 10, Sentences: []string{"すもももももももものうち"}}},
		result,
	); diff != "" {
		t.Errorf(diff)
	}
}

func newSentenceSplitterMock(ret []string) *sentenceSplitterMock {
	return &sentenceSplitterMock{
		splitArgs: []string{},
		splitRet:  ret,
	}
}

type sentenceSplitterMock struct {
	splitArgs []string
	splitRet  []string
}

func (s *sentenceSplitterMock) split(text string) ([]string, error) {
	s.splitArgs = append(s.splitArgs, text)
	return s.splitRet, nil
}

func newTokenizerMock(analyzeRet [][]string) *tokenizerMock {
	return &tokenizerMock{
		analyzeRet:  analyzeRet,
		analyzeArgs: [][]string{},
	}
}

type tokenizerMock struct {
	analyzeRet  [][]string
	analyzeArgs [][]string
}

func (t *tokenizerMock) analyze(text []string) [][]string {
	t.analyzeArgs = append(t.analyzeArgs, text)
	return t.analyzeRet
}

func newCharFilterMock(filterRet []string) *charFilterMock {
	return &charFilterMock{
		filterRet:  filterRet,
		fitlerArgs: [][]string{},
	}
}

type charFilterMock struct {
	filterRet  []string
	fitlerArgs [][]string
}

func (f *charFilterMock) filter(arg []string) []string {
	f.fitlerArgs = append(f.fitlerArgs, arg)
	return f.filterRet
}

func newWordFilterMock(filterRet [][]string) *wordFilterMock {
	return &wordFilterMock{
		filterRet:  filterRet,
		fitlerArgs: [][][]string{},
	}
}

type wordFilterMock struct {
	filterRet  [][]string
	fitlerArgs [][][]string
}

func (f *wordFilterMock) filter(arg [][]string) [][]string {
	f.fitlerArgs = append(f.fitlerArgs, arg)
	return f.filterRet
}

func newDbMock(searchRet []searchResult) *dbMock {
	return &dbMock{
		storeArgs:  []dbStoreArgs{},
		searchArgs: []dbSearchArgs{},
		searchRet:  searchRet,
	}
}

type dbMock struct {
	storeArgs  []dbStoreArgs
	searchArgs []dbSearchArgs
	searchRet  []searchResult
}

type dbStoreArgs struct {
	Token     string
	Positions []position
	Uri       string
}

type dbSearchArgs struct {
	Tokens []string
	Offset uint
	Count  uint
}

func (db *dbMock) store(token string, positions []position, uri string) error {
	db.storeArgs = append(db.storeArgs, dbStoreArgs{
		Token:     token,
		Positions: positions,
		Uri:       uri,
	})
	return nil
}

func (db *dbMock) search(tokens []string, offset, count uint) ([]searchResult, error) {
	db.searchArgs = append(db.searchArgs, dbSearchArgs{
		Tokens: tokens,
		Offset: offset,
		Count:  count,
	})
	return db.searchRet, nil
}
