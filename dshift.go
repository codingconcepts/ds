package main

import (
	"context"
	"database/sql"
	"dshift/internal/pkg/model"
	"dshift/internal/pkg/repo"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	version    string
	configPath string
)

func main() {
	log.SetFlags(0)

	rootCmd := &cobra.Command{
		Use:   "dshift",
		Short: "Shift data from one from database to another",
	}
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "absolute or relative path to the config file")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Print dshift version information",
			Run:   runVersion,
		},
		&cobra.Command{
			Use:   "insert",
			Short: "Insert data from one database into another",
			Run:   runInsert,
		},
		&cobra.Command{
			Use:   "update",
			Short: "Bring the target database up-to-date with the source database",
			Run:   runUpdate,
		},
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runVersion(cmd *cobra.Command, args []string) {
	log.Println(version)
}

func runInsert(cmd *cobra.Command, args []string) {
	config := loadConfig()

	sourceDB, err := sql.Open(config.Source.Driver, config.Source.URL)
	if err != nil {
		log.Fatalf("error connecting to source database: %v", err)
	}
	defer sourceDB.Close()

	targetDB, err := pgxpool.New(context.Background(), config.Target.URL)
	if err != nil {
		log.Fatalf("error connecting to target database: %v", err)
	}
	defer targetDB.Close()

	if err = repo.EnsureStateTable(targetDB, config.Target, false); err != nil {
		log.Fatalf("error ensuring state table: %v", err)
	}

	for _, sourceTable := range config.Source.Tables {
		targetTable, err := config.Target.GetTargetTable(sourceTable.SourceName)
		if err != nil {
			log.Fatalf("error getting target table: %v", err)
		}

		if err = repo.InsertTable(sourceDB, targetDB, sourceTable, targetTable); err != nil {
			log.Fatalf("error inserting %s -> %s: %v", sourceTable.Name, targetTable.Name, err)
		}
	}
}

func runUpdate(cmd *cobra.Command, args []string) {
	config := loadConfig()

	sourceDB, err := sql.Open(config.Source.Driver, config.Source.URL)
	if err != nil {
		log.Fatalf("error connecting to source database: %v", err)
	}
	defer sourceDB.Close()

	targetDB, err := pgxpool.New(context.Background(), config.Target.URL)
	if err != nil {
		log.Fatalf("error connecting to target database: %v", err)
	}
	defer targetDB.Close()

	if err = repo.EnsureStateTable(targetDB, config.Target, true); err != nil {
		log.Fatalf("error ensuring state table: %v", err)
	}

	for _, sourceTable := range config.Source.Tables {
		targetTable, err := config.Target.GetTargetTable(sourceTable.SourceName)
		if err != nil {
			log.Fatalf("error getting target table: %v", err)
		}

		if err = repo.UpdateTable(sourceDB, targetDB, sourceTable, targetTable); err != nil {
			log.Fatalf("error updating %s -> %s: %v", sourceTable.Name, targetTable.Name, err)
		}
	}
}

func loadConfig() model.Config {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("error opening config file: %v", err)
	}

	var c model.Config
	if err = yaml.NewDecoder(f).Decode(&c); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	return c
}
