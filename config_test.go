package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfig1(t *testing.T) {
	actual, _ := loadConfig("config1", []string{"test"})

	diff := cmp.Diff(
		config{
			Listen: "0.0.0.0:8000",
			Dsn:    "user:pass@tcp(127.0.0.1:3306)/searcher?charset=utf8&parseTime=True&loc=Local",
		},
		*actual,
	)
	if diff != "" {
		t.Errorf(diff)
	}
}

func TestLoadConfig2(t *testing.T) {
	actual, _ := loadConfig("config2", []string{"test"})

	if diff := cmp.Diff(
		config{
			Listen: "127.0.0.1:3000",
			Dsn:    "test:test@tcp(127.0.0.1:3306)/searcher?charset=utf8&parseTime=True&loc=Local",
		},
		*actual,
	); diff != "" {
		t.Errorf(diff)
	}
}
