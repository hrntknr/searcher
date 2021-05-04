package main

import (
	_ "embed"
	"encoding/json"

	"github.com/hrntknr/searcher/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed data/stopWords.json
var stopWordsData []byte

//go:embed data/mappingChar.json
var mappingCharData []byte

func main() {
	s, err := NewSearcher()
	if err != nil {
		panic(err)
	}
	if err := s.start(); err != nil {
		panic(err)
	}
}

func NewSearcher() (*Sercher, error) {
	config, err := loadConfig("config", []string{
		"/etc/searcher/",
		"$HOME/searcher/",
		".",
	})
	if err != nil {
		return nil, err
	}

	sentenceSplitter, err := newSentenceSplitter()
	if err != nil {
		return nil, err
	}

	tokenizer, err := newTokenizer()
	if err != nil {
		return nil, err
	}

	mappingChar := map[string]string{}
	if err := json.Unmarshal([]byte(mappingCharData), &mappingChar); err != nil {
		return nil, err
	}
	MappingCharFilter, err := newMappingCharFilter(mappingChar)
	if err != nil {
		return nil, err
	}
	lowercaseFilter, err := newLowercaseFilter()
	if err != nil {
		return nil, err
	}
	stopWords := []string{}
	if err := json.Unmarshal([]byte(stopWordsData), &stopWords); err != nil {
		return nil, err
	}
	stopWordFilter, err := newStopWordFilter(stopWords)
	if err != nil {
		return nil, err
	}
	stemmerFilter, err := newStemmerFilter()
	if err != nil {
		return nil, err
	}

	sql, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := sql.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(0)

	if err := sql.AutoMigrate(&types.Document{}, &types.Sentence{}, &types.Posting{}, &types.Token{}); err != nil {
		return nil, err
	}
	db, err := newDb(sql)
	if err != nil {
		return nil, err
	}

	service, err := newService(
		sentenceSplitter,
		tokenizer,
		[]CharFilter{MappingCharFilter},
		[]WordFilter{lowercaseFilter, stopWordFilter, stemmerFilter},
		db,
	)
	if err != nil {
		return nil, err
	}

	controller, err := newController(config, service)
	if err != nil {
		return nil, err
	}

	return &Sercher{
		controller: controller,
	}, nil
}

type Sercher struct {
	controller *controller
}

func (s *Sercher) start() error {
	return s.controller.start()
}
