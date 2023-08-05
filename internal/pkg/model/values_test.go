package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToParams(t *testing.T) {
	sourceValues := Values{
		[]any{"a", 1, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)},
		[]any{"b", 2, time.Date(2023, 2, 2, 2, 2, 2, 2, time.UTC)},
		[]any{"c", 3, time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)},
	}

	act, err := sourceValues.ToParams()
	assert.Nil(t, err)
	assert.Equal(t, "($1, $2, $3), ($4, $5, $6), ($7, $8, $9)", act)
}

func TestFlatten(t *testing.T) {
	sourceValues := Values{
		[]any{"a", 1, time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)},
		[]any{"b", 2, time.Date(2023, 2, 2, 2, 2, 2, 2, time.UTC)},
		[]any{"c", 3, time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)},
	}

	act := sourceValues.Flatten()
	exp := []any{
		"a",
		1,
		time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC),
		"b",
		2,
		time.Date(2023, 2, 2, 2, 2, 2, 2, time.UTC),
		"c",
		3,
		time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC),
	}
	assert.Equal(t, exp, act)
}
