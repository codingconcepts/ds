package repo

import "database/sql"

type mockRows struct {
	rows    [][]any
	columns []string
	index   int
}

func newMockRows(rows [][]any, columns []string) *mockRows {
	return &mockRows{
		rows:    rows,
		columns: columns,
	}
}

func (m *mockRows) Scan(dest ...interface{}) error {
	if m.index >= len(m.rows) {
		return sql.ErrNoRows
	}

	copy(dest, m.rows[m.index])

	m.index++
	return nil
}

func (m *mockRows) Next() bool {
	return m.index < len(m.rows)
}

func (m *mockRows) Columns() ([]string, error) {
	return m.columns, nil
}
