package repo

import (
	"ds/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	columns := []string{"a", "b", "c", "d"}

	rows := [][]any{
		{1, "A", "2023-01-01", true},
		{2, "B", "2023-01-02", false},
		{3, "C", "2023-01-03", true},
	}

	mockRows := newMockRows(rows, columns)

	table := model.Table{
		Columns: []model.Column{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
			{Name: "d"},
		},
	}

	exp := model.Values{
		[]any{1, "A", "2023-01-01", true},
		[]any{2, "B", "2023-01-02", false},
		[]any{3, "C", "2023-01-03", true},
	}

	act, err := scan(mockRows, table)
	assert.Nil(t, err)
	assert.Equal(t, exp, act)
}
