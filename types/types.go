package types

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Uri        string
	Time       time.Time
	TokenCount uint
}

type Sentence struct {
	gorm.Model
	DocumentID uint
	Index      uint
	Sentence   string
	TokenCount uint
	Postings   []*Posting `gorm:"many2many:posting_sentences"`
}

type Posting struct {
	gorm.Model
	TokenID    uint
	DocumentID uint
	Sentences  []*Sentence `gorm:"many2many:posting_sentences"`
}

type Token struct {
	gorm.Model
	Token string
}

type SearchResult struct {
	Uri       string
	Score     float64
	Sentences []string
}
