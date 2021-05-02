package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func init() {
	// ignore message
	gin.SetMode("release")
}

func TestControllerRegist(t *testing.T) {
	config, _ := loadConfig("config", []string{"test"})
	serviceMock := newServiceMock()
	controller, _ := newController(config, serviceMock)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/regist", bytes.NewBufferString("{\"uri\":\"test\",\"body\":\"すもももももももものうち\"}"))
	controller.router.ServeHTTP(w, req)

	if diff := cmp.Diff(
		200,
		w.Code,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[][]string{{"test", "すもももももももものうち"}},
		serviceMock.registArgs,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestControllerSearch(t *testing.T) {
	config, _ := loadConfig("config", []string{"test"})
	serviceMock := newServiceMock()
	controller, _ := newController(config, serviceMock)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/search?k=すもも&offset=11&count=12", nil)
	controller.router.ServeHTTP(w, req)

	if diff := cmp.Diff(
		200,
		w.Code,
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		`[{"Uri":"uri","Score":10,"Sentences":["すもももももももものうち"]}]`,
		string(w.Body.Bytes()),
	); diff != "" {
		t.Errorf(diff)
	}

	if diff := cmp.Diff(
		[]controllerSearchArgs{{Keyword: "すもも", Count: 12, Offset: 11}},
		serviceMock.searchArgs,
	); diff != "" {
		t.Errorf(diff)
	}
}

func newServiceMock() *serviceMock {
	return &serviceMock{
		registArgs: [][]string{},
		searchArgs: []controllerSearchArgs{},
	}
}

type serviceMock struct {
	registArgs [][]string
	searchArgs []controllerSearchArgs
}

func (s *serviceMock) regist(uri string, body string) error {
	s.registArgs = append(s.registArgs, []string{uri, body})
	return nil
}

func (s *serviceMock) search(keyword string, offset, count uint) ([]searchResult, error) {
	s.searchArgs = append(s.searchArgs, controllerSearchArgs{Keyword: keyword, Offset: offset, Count: count})
	return []searchResult{
		{
			Uri:       "uri",
			Score:     10,
			Sentences: []string{"すもももももももものうち"},
		},
	}, nil
}

type controllerSearchArgs struct {
	Keyword string
	Offset  uint
	Count   uint
}
