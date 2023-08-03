package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"shift/internal/pkg/model"
	"shift/internal/pkg/repo"

	"gopkg.in/yaml.v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var version string

func main() {
	log.SetFlags(0)

	configPath := flag.String("c", "", "absolute or relative path to the config file")
	versionFlag := flag.Bool("version", false, "display the current version number")
	flag.Parse()

	if *versionFlag {
		log.Println(version)
		return
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(*configPath)
	if err != nil {
		log.Fatalf("error opening config file: %v", err)
	}

	var c model.Config
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

	if err = repo.EnsureStateTable(targetDB, c.Target); err != nil {
		log.Fatalf("error ensuring state table: %v", err)
	}

	for _, sourceTable := range c.Source.Tables {
		targetTable, err := c.Target.GetTargetTable(sourceTable.SourceName)
		if err != nil {
			log.Fatalf("error getting target table: %v", err)
		}

		if err = repo.ShiftTable(sourceDB, targetDB, sourceTable, targetTable); err != nil {
			log.Fatalf("error shifting %s -> %s: %v", sourceTable.Name, targetTable.Name, err)
		}
	}
}
