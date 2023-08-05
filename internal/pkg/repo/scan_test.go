package repo

import (
	"ds/internal/pkg/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestScan(t *testing.T) {
	rows := sqlmock.NewRows(
		[]string{"a", "b", "c"}).
		AddRow(1, "A", "2023-01-01", true).
		AddRow(2, "B", "2023-01-02", false).
		AddRow(3, "C", "2023-01-03", true)

	table := model.Table{
		Columns: []model.Column{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
	}
	act, err := scan(rows)
}

func TestMapNArray(t *testing.T) {

}
