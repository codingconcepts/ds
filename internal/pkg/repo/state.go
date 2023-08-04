package repo

import (
	"context"
	"dshift/internal/pkg/model"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureStateTable creates the state table and initialises it with zeros for
// each of the migration tables.
func EnsureStateTable(targetDB *pgxpool.Pool, d model.Database, reset bool) error {
	// Create table if it doesn't exist.
	const tableStmt = `CREATE TABLE IF NOT EXISTS _shift_state (
		"table_name" STRING PRIMARY KEY,
		"current_offset" INT NOT NULL DEFAULT 0
	)`
	if _, err := targetDB.Exec(context.Background(), tableStmt); err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

	// Add tables if they don't exist.
	for _, table := range d.Tables {
		rowStmt := `INSERT INTO _shift_state (table_name) VALUES ($1)
								ON CONFLICT DO NOTHING`

		if _, err := targetDB.Exec(context.Background(), rowStmt, table.Name); err != nil {
			return fmt.Errorf("initialising table state: %w", err)
		}

		if !reset {
			continue
		}

		resetStmt := `UPDATE _shift_state SET current_offset = 0 WHERE true`
		if _, err := targetDB.Exec(context.Background(), resetStmt); err != nil {
			return fmt.Errorf("resetting table state: %w", err)
		}
	}

	return nil
}

// getShiftState returns the current_offset for a given table.
func getShiftState(targetDB *pgxpool.Pool, table string) (int, error) {
	const stmt = `SELECT current_offset FROM _shift_state WHERE table_name = $1`

	row := targetDB.QueryRow(context.Background(), stmt, table)

	var offset int
	if err := row.Scan(&offset); err != nil {
		return 0, fmt.Errorf("scanning row: %w", err)
	}

	return offset, nil
}

// setShiftState sets the current_offset for a given table.
func setShiftState(targetDB *pgxpool.Pool, table string, offset int) error {
	const stmt = `UPDATE _shift_state SET current_offset = $1 WHERE table_name = $2`

	if _, err := targetDB.Exec(context.Background(), stmt, offset, table); err != nil {
		return fmt.Errorf("updating offset: %w", err)
	}

	return nil
}
