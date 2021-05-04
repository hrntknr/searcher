package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hrntknr/searcher/mock"
	"github.com/hrntknr/searcher/types"
)

func init() {
	// ignore message
	gin.SetMode("release")
}

func TestControllerRegist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock.NewMockService(ctrl)
	gomock.InOrder(
		serviceMock.EXPECT().Regist("test", "すもももももももものうち"),
	)

	config, _ := loadConfig("config", []string{"test"})
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
}

func TestControllerSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock.NewMockService(ctrl)
	gomock.InOrder(
		serviceMock.EXPECT().Search("すもも", uint(11), uint(12)).Return(
			[]types.SearchResult{{
				Uri:       "uri",
				Score:     10,
				Sentences: []string{"すもももももももものうち"},
			}}, nil,
		),
	)

	config, _ := loadConfig("config", []string{"test"})
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
}
