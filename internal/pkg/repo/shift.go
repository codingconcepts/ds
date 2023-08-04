package repo

import (
	"context"
	"database/sql"
	"dshift/internal/pkg/model"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsertTable performs a bulk insert from the source database into the target database.
func InsertTable(sourceDB *sql.DB, targetDB *pgxpool.Pool, sourceTable, targetTable model.Table) error {
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

func UpdateTable(sourceDB *sql.DB, targetDB *pgxpool.Pool, sourceTable, targetTable model.Table) error {
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

		// Generate logical upsert statement.
		if stmt, err = targetTable.UpsertStatement(values); err != nil {
			return fmt.Errorf("generating upsert statement: %w", err)
		}

		if _, err = targetDB.Exec(context.Background(), stmt, values.Flatten()...); err != nil {
			return fmt.Errorf("upserting rows: %w", err)
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
