package psql

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPGStorage_GetAll(t *testing.T) {
	db, mock := NewMock()

	mock.ExpectQuery(`SELECT name, value FROM gauges WHERE type = ?`).WithArgs("gauge").WillReturnRows(
		sqlmock.NewRows([]string{"name", "value"}).AddRow("test_metric", 1.25))
	mock.ExpectQuery(`SELECT name, value FROM counters WHERE type = ?`).WithArgs("counter").WillReturnRows(
		sqlmock.NewRows([]string{"name", "value"}).AddRow("test_metric", 1.25))

	subject := PGStorage{
		conn: db,
	}

	resp, err := subject.GetAll(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp))
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestPGStorage_SetStrategy(t *testing.T) {
	pg1 := PGStorage{}
	pg1.SetStrategy(CountType)

	_, ok := pg1.strategy.(*Count)
	assert.True(t, ok)

	pg2 := PGStorage{}
	pg2.SetStrategy(GaugeType)

	_, ok = pg2.strategy.(*Gauge)
	assert.True(t, ok)
}
