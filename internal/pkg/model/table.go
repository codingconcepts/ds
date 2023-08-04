package model

import (
	"dshift/internal/pkg/io"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
)

// Table represents tables in the source and target databases.
type Table struct {
	Name string `yaml:"name"`

	// PrimaryKey is the column that uniquely identifies the row.
	PrimaryKey string `yaml:"primary_key"`

	// SourceName informs shift the origin table name, if target is different.
	SourceName string `yaml:"source_name"`

	// Filter informs shift that the origin table should be read from a particular point,
	// rather than in its entirety.
	Filter string `yaml:"filter"`

	// ReadLimit limits the number of rows to read from the source at any time.
	ReadLimit int `yaml:"read_limit"`

	// ReadDelay throttles reads from the source so neither database gets hammered.
	ReadDelay time.Duration `yaml:"read_delay"`
	Columns   []Column      `yaml:"columns"`
}

// SelectStatement returns a SELECT statement for a table's columns.
func (t Table) SelectStatement(offset int) string {
	columns := t.ColumnNames()

	return fmt.Sprintf(
		"SELECT %s FROM %s %s LIMIT %d OFFSET %d",
		strings.Join(columns, ","),
		t.Name,
		t.Filter,
		t.ReadLimit,
		offset,
	)
}

func (t Table) UpsertStatement(sourceValues Values) (string, error) {
	colums := t.ColumnNames()

	params, err := sourceValues.ToParams()
	if err != nil {
		return "", fmt.Errorf("creating params for upsert: %w", err)
	}

	fieldsForSet, err := t.fieldsForSetStatement()
	if err != nil {
		return "", fmt.Errorf("creating fields for set statement: %w", err)
	}

	return fmt.Sprintf(
		`INSERT INTO %s AS _shift_t (%s) VALUES %s
		 ON CONFLICT (%s) DO UPDATE
		 SET %s
		 WHERE _shift_t IS DISTINCT FROM EXCLUDED`,
		t.Name,
		strings.Join(colums, ","),
		params,
		t.PrimaryKey,
		fieldsForSet,
	), nil
}

func (t Table) fieldsForSetStatement() (string, error) {
	columns := t.ColumnNames()
	columns = lo.Reject(columns, func(col string, idx int) bool {
		return col == t.PrimaryKey
	})

	sb := io.NewErrWriter(strings.Builder{})
	for i, col := range columns {
		sb.WriteString(fmt.Sprintf("%s = EXCLUDED.%s", col, col))

		if i < len(columns)-1 {
			sb.WriteString(",")
		}
	}

	if err := sb.Err(); err != nil {
		return "", fmt.Errorf("generating set statement for update: %w", err)
	}

	return sb.String(), nil
}

// ColumnNames returns a slice of strings representing a table's column names.
func (t Table) ColumnNames() []string {
	return lo.Map(t.Columns, func(c Column, _ int) string {
		return c.Name
	})
}
