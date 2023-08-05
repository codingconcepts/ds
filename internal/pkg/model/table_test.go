package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelectStatement(t *testing.T) {
	cases := []struct {
		name   string
		table  Table
		offset int
		exp    string
	}{
		{
			name: "no filter or read limit",
			table: Table{
				Name: "test",
				Columns: []Column{
					{Name: "a"},
					{Name: "b"},
					{Name: "c"},
				}},
			offset: 0,
			exp:    `SELECT a, b, c FROM test   OFFSET 0`,
		},
		{
			name: "filter and read limit",
			table: Table{
				Name:      "test",
				Filter:    "WHERE col < '2023-01-01'",
				ReadLimit: 10,
				Columns: []Column{
					{Name: "a"},
					{Name: "b"},
					{Name: "c"},
				}},
			offset: 0,
			exp:    `SELECT a, b, c FROM test WHERE col < '2023-01-01' LIMIT 10 OFFSET 0`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			act := c.table.SelectStatement(c.offset)
			assert.Equal(t, c.exp, act)
		})
	}
}

func TestUpsertStatement(t *testing.T) {
	table := Table{
		Name:      "test",
		ReadLimit: 10,
		Columns: []Column{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		PrimaryKey: "id",
	}

	sourceValues := Values{
		[]any{"a", 1, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)},
		[]any{"b", 2, time.Date(2023, 2, 2, 2, 2, 2, 2, time.UTC)},
		[]any{"c", 3, time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)},
	}

	act, err := table.UpsertStatement(sourceValues)
	assert.Nil(t, err)

	exp := "INSERT INTO test AS _shift_t (a, b, c) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)\n\t\t ON CONFLICT (id) DO UPDATE\n\t\t SET a = EXCLUDED.a, b = EXCLUDED.b, c = EXCLUDED.c\n\t\t WHERE _shift_t IS DISTINCT FROM EXCLUDED"
	assert.Equal(t, exp, act)
}

func TestColumnNames(t *testing.T) {

}
