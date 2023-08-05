package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTargetTable(t *testing.T) {
	cases := []struct {
		name         string
		database     Database
		source       string
		expTableName string
		expErr       error
	}{
		{
			name: "target table found from source name",
			database: Database{
				Tables: []Table{
					{
						Name:       "a",
						SourceName: "b",
					},
				},
			},
			source:       "b",
			expTableName: "a",
		},
		{
			name: "target table found from name",
			database: Database{
				Tables: []Table{
					{
						Name:       "a",
						SourceName: "b",
					},
				},
			},
			source:       "a",
			expTableName: "a",
		},
		{
			name: "target table not found",
			database: Database{
				Tables: []Table{
					{
						Name:       "a",
						SourceName: "b",
					},
				},
			},
			source: "c",
			expErr: fmt.Errorf("missing target for c; ensure table names match, target has a source_name"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			act, err := c.database.GetTargetTable(c.source)
			assert.Equal(t, c.expErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.expTableName, act.Name)
		})
	}
}
