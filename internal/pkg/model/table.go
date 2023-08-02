package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
)

// Table represents tables in the source and target databases.
type Table struct {
	Name string `yaml:"name"`

	// SourceName informs shift the origin table name, if target is different.
	SourceName string `yaml:"source_name"`

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
		"SELECT %s FROM %s LIMIT %d OFFSET %d",
		strings.Join(columns, ", "),
		t.Name,
		t.ReadLimit,
		offset,
	)
}

// ColumnNames returns a slice of strings representing a table's column names.
func (t Table) ColumnNames() []string {
	return lo.Map(t.Columns, func(c Column, _ int) string {
		return c.Name
	})
}
