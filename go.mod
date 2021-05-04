module github.com/hrntknr/searcher

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/gin-gonic/gin v1.7.1
	github.com/golang/mock v1.5.0
	github.com/google/go-cmp v0.5.5
	github.com/hrntknr/searcher/mock v0.0.0-00010101000000-000000000000
	github.com/hrntknr/searcher/types v0.0.0-00010101000000-000000000000
	github.com/ikawaha/kagome-dict/ipa v1.0.2
	github.com/ikawaha/kagome/v2 v2.4.4
	github.com/kljensen/snowball v0.6.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	gorm.io/driver/mysql v1.0.6
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.9
)

replace github.com/hrntknr/searcher/mock => ./mock

replace github.com/hrntknr/searcher/types => ./types
