package repo

import (
	"ds/internal/pkg/model"
	"fmt"
)

type rowScanner interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}

// scan a row collection for a given table into a multi-dimensional array.
func scan(rs rowScanner, t model.Table) (model.Values, error) {
	fields, err := rs.Columns()
	if err != nil {
		return nil, fmt.Errorf("listing columns: %w", err)
	}

	var rows []map[string]any
	for rs.Next() {
		scans := make([]any, len(fields))
		row := make(map[string]any)

		for i := range scans {
			scans[i] = &scans[i]
		}

		if err = rs.Scan(scans...); err != nil {
			return nil, fmt.Errorf("scaning values: %w", err)
		}

		for i, v := range scans {
			if v != nil {
				row[fields[i]] = v
			}
		}
		rows = append(rows, row)
	}

	return mapToNArray(rows, t), nil
}

func mapToNArray(m []map[string]any, t model.Table) model.Values {
	values := model.Values{}

	for _, row := range m {
		columns := make([]any, len(t.Columns))
		for i, col := range t.Columns {
			columns[i] = row[col.Name]
		}
		values = append(values, columns)
	}

	return values
}
