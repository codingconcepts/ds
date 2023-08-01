package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config file: %v", err)
	}

	var c config
	if err = yaml.NewDecoder(f).Decode(&c); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	sourceDB, err := sql.Open(c.Source.Driver, c.Source.URL)
	if err != nil {
		log.Fatalf("error connecting to source database: %v", err)
	}
	defer sourceDB.Close()

	targetDB, err := pgxpool.New(context.Background(), c.Target.URL)
	if err != nil {
		log.Fatalf("error connecting to target database: %v", err)
	}
	defer targetDB.Close()

	if err = ensureStateTable(targetDB, c.Target); err != nil {
		log.Fatalf("error ensuring state table: %v", err)
	}

	for _, table := range c.Source.Tables {
		if err = shiftTable(sourceDB, targetDB, table); err != nil {
			log.Fatalf("error shifting table: %s: %v", table.Name, err)
		}
	}
}

func shiftTable(sourceDB *sql.DB, targetDB *pgxpool.Pool, t table) error {
	for {
		// Fetch current offset.
		offset, err := getShiftState(targetDB, t.Name)
		if err != nil {
			return fmt.Errorf("fetching current offset: %w", err)
		}

		// Read from input.
		stmt := t.selectStatement(offset)
		rows, err := sourceDB.Query(stmt)
		if err != nil {
			return fmt.Errorf("querying rows: %w", err)
		}

		values, err := scan(rows, t)
		if err != nil {
			return fmt.Errorf("scanning rows: %w", err)
		}

		if len(values) == 0 {
			return nil
		}

		// Write to output.
		if _, err = targetDB.CopyFrom(context.Background(), pgx.Identifier{t.Name}, t.columnNames(), pgx.CopyFromRows(values)); err != nil {
			return fmt.Errorf("inserting rows: %w", err)
		}

		// Set current offset.
		if err = setShiftState(targetDB, t.Name, offset+len(values)); err != nil {
			return fmt.Errorf("setting current offset: %w", err)
		}

		// Exit loop if we've read less than the read_limit.
		if len(values) < t.ReadLimit {
			return nil
		}
	}
}

func getShiftState(targetDB *pgxpool.Pool, table string) (int, error) {
	const stmt = `SELECT current_offset FROM _shift_state WHERE table_name = $1`

	row := targetDB.QueryRow(context.Background(), stmt, table)

	var offset int
	if err := row.Scan(&offset); err != nil {
		return 0, fmt.Errorf("scanning row: %w", err)
	}

	return offset, nil
}

func setShiftState(targetDB *pgxpool.Pool, table string, offset int) error {
	const stmt = `UPDATE _shift_state SET current_offset = $1 WHERE table_name = $2`

	if _, err := targetDB.Exec(context.Background(), stmt, offset, table); err != nil {
		return fmt.Errorf("updating offset: %w", err)
	}

	return nil
}

func ensureStateTable(targetDB *pgxpool.Pool, d database) error {
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
	}

	return nil
}

func scan(list *sql.Rows, t table) ([][]any, error) {
	fields, err := list.Columns()
	if err != nil {
		return nil, fmt.Errorf("listing columns: %w", err)
	}

	var rows []map[string]any
	for list.Next() {
		scans := make([]any, len(fields))
		row := make(map[string]any)

		for i := range scans {
			scans[i] = &scans[i]
		}

		if err = list.Scan(scans...); err != nil {
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

func mapToNArray(m []map[string]any, t table) [][]any {
	array := [][]any{}

	for _, row := range m {
		columns := make([]any, len(t.Columns))
		for i, col := range t.Columns {
			columns[i] = row[col.Name]
		}
		array = append(array, columns)
	}

	return array
}

type config struct {
	Source database `yaml:"source"`
	Target database `yaml:"target"`
}

type database struct {
	Driver string  `yaml:"driver"`
	URL    string  `yaml:"url"`
	Tables []table `yaml:"tables"`
}

type table struct {
	Name      string   `yaml:"name"`
	ReadLimit int      `yaml:"read_limit"`
	Columns   []column `yaml:"columns"`
}

func (t table) selectStatement(offset int) string {
	columns := t.columnNames()

	return fmt.Sprintf(
		"SELECT %s FROM %s LIMIT %d OFFSET %d",
		strings.Join(columns, ", "),
		t.Name,
		t.ReadLimit,
		offset,
	)
}

func (t table) columnNames() []string {
	return lo.Map(t.Columns, func(c column, _ int) string {
		return c.Name
	})
}

type column struct {
	Name string `yaml:"name"`
}
