package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hrntknr/searcher/types"
	"github.com/stretchr/testify/assert"
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

func TestCountDocument(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "documents" WHERE "documents"."deleted_at" IS NULL`,
	)).WillReturnRows(
		sqlmock.NewRows([]string{"count(1)"}).
			AddRow(100),
	)

	count, err := db.CountDocument()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, count, uint(100))
}

func TestCountTermInDocument(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "documents" WHERE id = $1 AND "documents"."deleted_at" IS NULL ORDER BY "documents"."id" LIMIT 1`,
	)).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"token_count"}).
			AddRow(10),
	)

	count, err := db.CountTermInDocument(1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, count, uint(10))
}

func TestDocumentFromUri(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "documents" WHERE uri = $1 AND "documents"."deleted_at" IS NULL ORDER BY "documents"."id" LIMIT 1`,
	)).WithArgs("uri").WillReturnRows(
		sqlmock.NewRows([]string{"id", "uri", "time", "token_count"}).
			AddRow(1, "uri", time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC), 100),
	)

	document, err := db.DocumentFromUri("uri")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Document{
			Model: gorm.Model{
				ID: 1,
			},
			Uri:        "uri",
			Time:       time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
			TokenCount: 100,
		},
		document,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestDocumentFromID(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "documents" WHERE id = $1 AND "documents"."deleted_at" IS NULL ORDER BY "documents"."id" LIMIT 1`,
	)).WithArgs(10).WillReturnRows(
		sqlmock.NewRows([]string{"id", "uri", "time", "token_count"}).
			AddRow(10, "uri", time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC), 100),
	)

	document, err := db.DocumentFromID(10)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Document{
			Model: gorm.Model{
				ID: 10,
			},
			Uri:        "uri",
			Time:       time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
			TokenCount: 100,
		},
		document,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestCreateDcoument(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "documents" ("created_at","updated_at","deleted_at","uri","time","token_count") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
	)).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		"uri",
		time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
		100,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(10),
	)
	mock.ExpectCommit()

	document, err := db.CreateDcoument(&types.Document{
		Uri:        "uri",
		Time:       time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
		TokenCount: 100,
	})
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Document{
			Model: gorm.Model{
				ID: 10,
			},
			Uri:        "uri",
			Time:       time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
			TokenCount: 100,
		},
		document,
		cmpopts.IgnoreFields(*document, "Model.CreatedAt"),
		cmpopts.IgnoreFields(*document, "Model.UpdatedAt"),
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestTokenFromString(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tokens" WHERE token = $1 AND "tokens"."deleted_at" IS NULL ORDER BY "tokens"."id" LIMIT 1`,
	)).WithArgs("token").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token"}).
			AddRow(10, "token"),
	)

	token, err := db.TokenFromString("token")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Token{
			Model: gorm.Model{
				ID: 10,
			},
			Token: "token",
		},
		token,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestTokenFromID(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tokens" WHERE id = $1 AND "tokens"."deleted_at" IS NULL ORDER BY "tokens"."id" LIMIT 1`,
	)).WithArgs(10).WillReturnRows(
		sqlmock.NewRows([]string{"id", "token"}).
			AddRow(10, "token"),
	)

	token, err := db.TokenFromID(10)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Token{
			Model: gorm.Model{
				ID: 10,
			},
			Token: "token",
		},
		token,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestCreateToken(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tokens" ("created_at","updated_at","deleted_at","token") VALUES ($1,$2,$3,$4) RETURNING "id"`,
	)).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		"token",
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(11),
	)
	mock.ExpectCommit()

	token, err := db.CreateToken(&types.Token{
		Token: "token",
	})
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Token{
			Model: gorm.Model{
				ID: 11,
			},
			Token: "token",
		},
		token,
		cmpopts.IgnoreFields(*token, "Model.CreatedAt"),
		cmpopts.IgnoreFields(*token, "Model.UpdatedAt"),
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestPostingList(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "postings" WHERE token_id = $1 AND "postings"."deleted_at" IS NULL`,
	)).WithArgs(10).WillReturnRows(
		sqlmock.NewRows([]string{"id", "token_id", "document_id"}).
			AddRow(10, 10, 1),
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "posting_sentences" WHERE "posting_sentences"."posting_id" = $1`,
	)).WithArgs(10).WillReturnRows(
		sqlmock.NewRows([]string{"id", "posting_id", "sentence_id"}).
			AddRow(10, 10, 12).
			AddRow(10, 10, 13),
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sentences" WHERE "sentences"."id" IN ($1,$2) AND "sentences"."deleted_at" IS NULL`,
	)).WithArgs(12, 13).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).
			AddRow(12).
			AddRow(13),
	)

	ps, err := db.PostingList(10)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		[]*types.Posting{{
			Model: gorm.Model{
				ID: 10,
			},
			TokenID:    10,
			DocumentID: 1,
			Sentences: []*types.Sentence{{
				Model: gorm.Model{
					ID: 12,
				},
			}, {
				Model: gorm.Model{
					ID: 13,
				},
			}},
		}},
		ps,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestCreatePosting(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "postings" ("created_at","updated_at","deleted_at","token_id","document_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
	)).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		1,
		2,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(10),
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "sentences" ("created_at","updated_at","deleted_at","document_id","index","sentence","token_count") VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT DO NOTHING RETURNING "id"`,
	)).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		2,
		0,
		"test",
		4,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(11),
	)
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "posting_sentences" ("posting_id","sentence_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`,
	)).WithArgs(
		10,
		11,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()

	posting, err := db.CreatePosting(&types.Posting{
		TokenID:    1,
		DocumentID: 2,
		Sentences: []*types.Sentence{{
			DocumentID: 2,
			Index:      0,
			Sentence:   "test",
			TokenCount: 4,
		}},
	})
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Posting{
			Model: gorm.Model{
				ID: 10,
			},
			TokenID:    1,
			DocumentID: 2,
			Sentences: []*types.Sentence{{
				Model: gorm.Model{
					ID: 11,
				},
				DocumentID: 2,
				Index:      0,
				Sentence:   "test",
				TokenCount: 4,
			}},
		},
		posting,
		cmpopts.IgnoreFields(*posting, "Model.CreatedAt"),
		cmpopts.IgnoreFields(*posting, "Model.UpdatedAt"),
		cmpopts.IgnoreFields(*posting, "Sentences", "Model.CreatedAt"),
		cmpopts.IgnoreFields(*posting, "Sentences", "Model.UpdatedAt"),
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestSentenceMultiFromID(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sentences" WHERE "sentences"."id" IN ($1,$2,$3) AND "sentences"."deleted_at" IS NULL`,
	)).WithArgs(
		1, 2, 3,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).
			AddRow(1).
			AddRow(2).
			AddRow(3),
	)

	sentences, err := db.SentenceMultiFromID([]uint{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		[]*types.Sentence{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 2}},
			{Model: gorm.Model{ID: 3}},
		},
		sentences,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestCreateSentence(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "sentences" ("created_at","updated_at","deleted_at","document_id","index","sentence","token_count") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`,
	)).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nil,
		1,
		0,
		"test",
		100,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(10),
	)
	mock.ExpectCommit()

	sentence, err := db.CreateSentence(&types.Sentence{
		DocumentID: 1,
		Index:      0,
		Sentence:   "test",
		TokenCount: 100,
	})
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		&types.Sentence{
			Model: gorm.Model{
				ID: 10,
			},
			DocumentID: 1,
			Index:      0,
			Sentence:   "test",
			TokenCount: 100,
		},
		sentence,
		cmpopts.IgnoreFields(*sentence, "Model.CreatedAt"),
		cmpopts.IgnoreFields(*sentence, "Model.UpdatedAt"),
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestDeleteSentenceFromDocumentID(t *testing.T) {
	gdb, mock, _ := getDBMock()
	db, _ := newDb(gdb)
	mock.ExpectBegin()
	// mock.ExpectExec(regexp.QuoteMeta(
	// 	`DELETE FROM "posting_sentences" WHERE "posting_sentences"."sentence_id" IN (NULL)`,
	// )).WithArgs().WillReturnResult(
	// 	sqlmock.NewResult(1, 1),
	// )
	// mock.ExpectExec(regexp.QuoteMeta(
	// 	`UPDATE "sentences" SET "deleted_at"=$1 WHERE document_id = $2 AND "sentences"."deleted_at" IS NULL`,
	// )).WithArgs(
	// 	sqlmock.AnyArg(),
	// 	10,
	// ).WillReturnResult(
	// 	sqlmock.NewResult(1, 1),
	// )
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sentences" WHERE document_id = $1 AND "sentences"."deleted_at" IS NULL`,
	)).WithArgs(10).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(11),
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "posting_sentences" WHERE "posting_sentences"."sentence_id" = $1`,
	)).WithArgs(11).WillReturnRows(
		sqlmock.NewRows([]string{"id", "sentence_id", "posting_id"}).AddRow(12, 11, 13),
	)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "postings" WHERE "postings"."id" = $1 AND "postings"."deleted_at" IS NULL`,
	)).WithArgs(13).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(13),
	)
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "posting_sentences" WHERE "posting_sentences"."sentence_id" = $1`,
	)).WithArgs(11).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "sentences" SET "deleted_at"=$1 WHERE "sentences"."id" = $2 AND "sentences"."deleted_at" IS NULL`,
	)).WithArgs(sqlmock.AnyArg(), 11).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "postings" SET "deleted_at"=$1 WHERE "postings"."id" = $2 AND "postings"."deleted_at" IS NULL`,
	)).WithArgs(sqlmock.AnyArg(), 13).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()
	if err := db.DeleteSentenceFromDocumentID(10); err != nil {
		t.Error(err)
	}
}
