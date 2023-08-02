package repo

import (
	"context"
	"database/sql"
	"fmt"
	"shift/internal/pkg/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ShiftTable from the source database to the target database.
func ShiftTable(sourceDB *sql.DB, targetDB *pgxpool.Pool, sourceTable, targetTable model.Table) error {
	for {
		// Fetch current offset.
		offset, err := getShiftState(targetDB, sourceTable.Name)
		if err != nil {
			return fmt.Errorf("fetching current offset: %w", err)
		}

		// Read from input.
		stmt := sourceTable.SelectStatement(offset)
		rows, err := sourceDB.Query(stmt)
		if err != nil {
			return fmt.Errorf("querying rows: %w", err)
		}

		values, err := scan(rows, sourceTable)
		if err != nil {
			return fmt.Errorf("scanning rows: %w", err)
		}

		if len(values) == 0 {
			return nil
		}

		// Write to output.
		if _, err = targetDB.CopyFrom(context.Background(), pgx.Identifier{sourceTable.Name}, targetTable.ColumnNames(), pgx.CopyFromRows(values)); err != nil {
			return fmt.Errorf("inserting rows: %w", err)
		}

		// Set current offset.
		if err = setShiftState(targetDB, sourceTable.Name, offset+len(values)); err != nil {
			return fmt.Errorf("setting current offset: %w", err)
		}

		// Exit loop if we've read less than the read_limit.
		if len(values) < sourceTable.ReadLimit {
			return nil
		}

		if sourceTable.ReadDelay > 0 {
			time.Sleep(sourceTable.ReadDelay)
		}
	}
}
