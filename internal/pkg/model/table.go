package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// Table represents tables in the source and target databases.
type Table struct {
	Name       string   `yaml:"name"`
	SourceName string   `yaml:"source_name"`
	ReadLimit  int      `yaml:"read_limit"`
	Columns    []Column `yaml:"columns"`
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
