package model

import (
	"fmt"

	"github.com/samber/lo"
)

// Database represents the source and target database value in the config file.
type Database struct {
	Driver string  `yaml:"driver"`
	URL    string  `yaml:"url"`
	Tables []Table `yaml:"tables"`
}

// GetTargetTable returns the target table for a given source table name.
func (d Database) GetTargetTable(source string) (Table, error) {
	targetTable, ok := lo.Find(d.Tables, func(t Table) bool {
		return t.SourceName == source || t.Name == source
	})

	if !ok {
		return Table{}, fmt.Errorf("missing target for %s; ensure table names match, target has a source_name", source)
	}

	return targetTable, nil
}
