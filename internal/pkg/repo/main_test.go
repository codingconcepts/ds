package repo

import (
	"context"
	"database/sql"
	"ds/internal/pkg/model"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	integrationTests bool
	source           *sql.DB
	target           *pgxpool.Pool
)

func TestMain(m *testing.M) {
	integrationTestsRaw, ok := os.LookupEnv("INTEGRATION_TESTS")
	if ok {
		var err error
		if integrationTests, err = strconv.ParseBool(integrationTestsRaw); err != nil {
			log.Fatalf("error parsing INTEGRATION_TESTS env var: %v", err)
		}
	}

	if integrationTests {
		setupDatabase()
	}

	exitCode := m.Run()

	if integrationTests {
		teardownDatabase()
	}

	os.Exit(exitCode)
}

func setupDatabase() {
	var err error

	sourceURL, ok := os.LookupEnv("SOURCE_URL")
	if !ok {
		log.Fatal("SOURCE_URL env var not found")
	}

	targetURL, ok := os.LookupEnv("TARGET_URL")
	if !ok {
		log.Fatal("TARGET_URL env var not found")
	}

	if source, err = sql.Open("pgx", sourceURL); err != nil {
		log.Fatalf("error connecting to source database: %v", err)
	}
	if err = source.Ping(); err != nil {
		log.Fatalf("error pinging source database: %v", err)
	}

	if target, err = pgxpool.New(context.Background(), targetURL); err != nil {
		log.Fatalf("error connecting to target database: %v", err)
	}
	if err = target.Ping(context.Background()); err != nil {
		log.Fatalf("error pinging target database: %v", err)
	}

	createSourceDatabase()
	createTargetDatabase()

	targetDatabase := model.Database{
		Tables: []model.Table{
			{Name: "person"},
		},
	}

	if err = EnsureStateTable(target, targetDatabase, true); err != nil {
		log.Fatalf("error ensuring database: %v", err)
	}
}

func createSourceDatabase() {
	createStmt := `
		CREATE TABLE person (
			id UUID PRIMARY KEY,
			full_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`

	if _, err := source.Exec(createStmt); err != nil {
		log.Fatalf("creating db: %v", err)
	}

	insertStmt := `
		INSERT INTO person (id, full_name, created_at) VALUES
			('af57040a-f393-45a1-aa71-828ed9b20ca8', 'a a', '2023-01-01T01:01:01Z'),
			('bc229cee-5387-4c36-b83e-7a46613071de', 'b b', '2023-01-01T01:01:02Z'),
			('ccb45142-cba0-4f97-9179-33c7e3d51e92', 'c c', '2023-01-01T01:01:03Z'),
			('dfbad599-b5b6-4c59-8859-b19db1dd2ac0', 'd d', '2023-01-01T01:01:04Z'),
			('eba7ea84-e57b-4816-8806-2faae31c2830', 'e e', '2023-01-01T01:01:05Z')`

	if _, err := source.Exec(insertStmt); err != nil {
		log.Fatalf("seeding source db: %v", err)
	}
}

func createTargetDatabase() {
	stmt := `
		CREATE TABLE person (
			id UUID PRIMARY KEY,
			full_name STRING NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`

	if _, err := target.Exec(context.Background(), stmt); err != nil {
		log.Fatalf("creating db: %v", err)
	}
}

func dropSourceDatabase() {
	stmt := `DROP TABLE person`

	if _, err := source.Exec(stmt); err != nil {
		log.Fatalf("dropping source db: %v", err)
	}
}

func dropTargetDatabase() {
	stmt := `DROP TABLE person`

	if _, err := target.Exec(context.Background(), stmt); err != nil {
		log.Fatalf("dropping source db: %v", err)
	}
}

func teardownDatabase() {
	dropSourceDatabase()
	dropTargetDatabase()

	source.Close()
	target.Close()
}
