package main

import (
	_ "embed"
	"encoding/json"

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
	mappingCharFilter, err := newMappingCharFilter(mappingChar)
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
	if err := sql.AutoMigrate(&invertedIndex{}, &posting{}, &sentence{}, &tokenPosition{}); err != nil {
		return nil, err
	}
	db, err := newDb(sql)
	if err != nil {
		return nil, err
	}

	service, err := newService(
		sentenceSplitter,
		tokenizer,
		[]charFilter{mappingCharFilter},
		[]wordFilter{lowercaseFilter, stopWordFilter, stemmerFilter},
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
