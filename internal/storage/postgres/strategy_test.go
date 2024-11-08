package psql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCount_Get(t *testing.T) {
	db, mock := NewMock()

	mock.ExpectQuery(`SELECT name, type, value FROM counters WHERE name = ?`).WithArgs("testmetric").WillReturnRows(
		sqlmock.NewRows([]string{"testmetric", "counter", "1"}))

	c := Count{
		conn: db,
	}
	_, err := c.Get(context.Background(), "counter", "testmetric")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestGauge_Get(t *testing.T) {
	db, mock := NewMock()

	mock.ExpectQuery(`SELECT name, type, value FROM gauges WHERE name = ?`).WithArgs("testmetric").WillReturnRows(
		sqlmock.NewRows([]string{"testmetric", "gauges", "1"}))

	g := Gauge{
		conn: db,
	}
	_, err := g.Get(context.Background(), "gauges", "testmetric")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
