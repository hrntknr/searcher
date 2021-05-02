package main

import (
	"math"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func TestDbStoreUpdate(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "inverted_indices" WHERE token = $1 AND "inverted_indices"."deleted_at" IS NULL ORDER BY "inverted_indices"."id" LIMIT 1`,
	)).WillReturnRows(sqlmock.NewRows([]string{"Token", "ID"}))
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "inverted_indices" ("created_at","updated_at","deleted_at","token") VALUES ($1,$2,$3,$4) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "token").WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "postings" WHERE (uri = $1 AND inverted_index_id = $2) AND "postings"."deleted_at" IS NULL ORDER BY "postings"."id" LIMIT 1`,
	)).WithArgs("uri", 1).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(2))
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "sentences" SET "deleted_at"=$1 WHERE posting_id = $2 AND "sentences"."deleted_at" IS NULL`,
	)).WithArgs(sqlmock.AnyArg(), 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "token_positions" SET "deleted_at"=$1 WHERE posting_id = $2 AND "token_positions"."deleted_at" IS NULL`,
	)).WithArgs(sqlmock.AnyArg(), 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "sentences" ("created_at","updated_at","deleted_at","posting_id","body") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 2, "token dayo~~").WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(3))
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "token_positions" ("created_at","updated_at","deleted_at","sentence_id","sentence_position","posting_id","posting_position") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 3, 0, 2, 0).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
	mock.ExpectCommit()

	if err := db.store("token", []position{{0, 0, "token dayo~~"}}, "uri"); err != nil {
		t.Error(err)
	}
}

func TestDbStoreCreate(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "inverted_indices" WHERE token = $1 AND "inverted_indices"."deleted_at" IS NULL ORDER BY "inverted_indices"."id" LIMIT 1`,
	)).WillReturnRows(sqlmock.NewRows([]string{"Token", "ID"})).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "inverted_indices" ("created_at","updated_at","deleted_at","token") VALUES ($1,$2,$3,$4) RETURNING "id"`,
	)).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "postings" WHERE (uri = $1 AND inverted_index_id = $2) AND "postings"."deleted_at" IS NULL ORDER BY "postings"."id" LIMIT 1`,
	)).WithArgs("uri", 1).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "postings" ("created_at","updated_at","deleted_at","inverted_index_id","uri","token_count") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1, "uri", 1).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(2))
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "sentences" ("created_at","updated_at","deleted_at","posting_id","body") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 2, "token dayo~~").WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(3))
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "token_positions" ("created_at","updated_at","deleted_at","sentence_id","sentence_position","posting_id","posting_position") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`,
	)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 3, 0, 2, 0).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
	mock.ExpectCommit()

	if err := db.store("token", []position{{0, 0, "token dayo~~"}}, "uri"); err != nil {
		t.Error(err)
	}
}

func TestDbSearch(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	tokens := []string{"word1", "word2", "word3"}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "postings" WHERE "postings"."deleted_at" IS NULL`,
	)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1000))
	mock.MatchExpectationsInOrder(true)
	for i, token := range tokens {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "inverted_indices" WHERE token = $1 AND "inverted_indices"."deleted_at" IS NULL ORDER BY "inverted_indices"."id" LIMIT 1`,
		)).WithArgs(token).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(30))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "postings" WHERE inverted_index_id = $1 AND "postings"."deleted_at" IS NULL`,
		)).WithArgs(30).WillReturnRows(sqlmock.NewRows([]string{"ID", "Uri", "TokenCount"}).AddRow(31, "uri1", 6).AddRow(32, "uri2", 9))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "sentences" WHERE posting_id = $1 AND "sentences"."deleted_at" IS NULL`,
		)).WithArgs(31).WillReturnRows(sqlmock.NewRows([]string{"ID", "Body"}).AddRow(40, token+" "+token))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "sentences" WHERE posting_id = $1 AND "sentences"."deleted_at" IS NULL`,
		)).WithArgs(32).WillReturnRows(sqlmock.NewRows([]string{"ID", "Body"}).AddRow(41, token+" "+token+" "+"test"))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE sentence_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(40).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*2+0, 33, 0).AddRow(35, i*2+1, 33, 1))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE sentence_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(41).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*2+0, 34, 0).AddRow(35, i*2+1, 34, 1))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE sentence_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(42).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*2+0, 34, 0).AddRow(35, i*2+1, 34, 1))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE sentence_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(43).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*2+0, 34, 0).AddRow(35, i*2+1, 34, 1))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE posting_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(31).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*2+0, 33, 0).AddRow(35, i*2+1, 33, 1))
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "token_positions" WHERE posting_id = $1 AND "token_positions"."deleted_at" IS NULL`,
		)).WithArgs(32).WillReturnRows(sqlmock.NewRows([]string{"ID", "PostingPosition", "SentenceID", "SentencePosition"}).AddRow(35, i*3+0, 34, 0).AddRow(35, i*3+1, 34, 1))
	}
	mock.MatchExpectationsInOrder(false)

	result, err := db.search(tokens, 0, 10)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		[]searchResult{
			{
				Uri:   "uri1",
				Score: math.Pow((float64(2)/float64(6))*math.Log(float64(1000)/float64(3)), 3),
				Sentences: []string{
					"word1 word1",
					"word2 word2",
					"word3 word3",
				},
			},
			{
				Uri:   "uri2",
				Score: math.Pow((float64(2)/float64(9))*math.Log(float64(1000)/float64(3)), 3),
				Sentences: []string{
					"word1 word1 test",
					"word2 word2 test",
					"word3 word3 test",
				},
			},
		},
		result,
	); diff != "" {
		t.Errorf(diff)
	}
}
