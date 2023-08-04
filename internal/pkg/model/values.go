package model

import (
	"dshift/internal/pkg/io"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// Values wraps a multi-dimensional array of any, allowing behaviour
// to be attached.
type Values [][]any

// ToParams returns Values as a comma-separated collection of bracketed values,
// such that Values = [][]{"a", 1}, {"b", 2} would return:
//
// ($1,$2),($3,$4)
func (v Values) ToParams() (string, error) {
	sb := io.NewErrWriter(strings.Builder{})

	argIdx := 0
	for i, row := range v {
		sb.WriteString("(")

		for j := range row {
			argIdx += 1
			sb.WriteString("$%d", argIdx)
			if j < len(row)-1 {
				sb.WriteString(",")
			}
		}

		sb.WriteString(")")

		if i < len(v)-1 {
			sb.WriteString(",")
		}
	}

	if err := sb.Err(); err != nil {
		return "", fmt.Errorf("writing params: %w", err)
	}

	return sb.String(), nil
}

// Flatten returns a flat version of Values, such that Values = [][]{"a", 1}, {"b", 2} would return:
//
// []any{"a", 1, "b", 2}
func (v Values) Flatten() []any {
	return lo.Flatten(v)
}
