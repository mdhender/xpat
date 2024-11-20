// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"github.com/mdhender/xpat/stdlib"
	"github.com/mdhender/xpat/stores/sqlite"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var (
	argsDb struct {
		force bool // if true, overwrite existing database
		paths struct {
			database string // path to the database file
		}
	}

	cmdDb = &cobra.Command{
		Use:   "db",
		Short: "Database management commands",
	}

	cmdDbCreate = &cobra.Command{
		Use:   "create",
		Short: "Create new database or database objects",
	}

	argsDbCreateDatabase struct {
		force bool // if true, overwrite existing database
	}

	cmdDbCreateDatabase = &cobra.Command{
		Use:   "database",
		Short: "create and initialize a new database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if argsDb.paths.database == "" {
				return fmt.Errorf("database: path is required\n")
			} else if path, err := filepath.Abs(argsDb.paths.database); err != nil {
				return fmt.Errorf("database: %v\n", err)
			} else {
				argsDb.paths.database = path
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("db: create: database  %s\n", argsDb.paths.database)

			// it is an error if the database already exists unless force is true.
			// in that case, we remove the database so that we can create it again.
			if !argsDbCreateDatabase.force {
				if ok, err := stdlib.IsFileExists(argsDb.paths.database); err != nil {
					log.Fatalf("db: %v\n", err)
				} else if ok {
					log.Printf("db: create: removing %s\n", argsDb.paths.database)
					if err := os.Remove(argsDb.paths.database); err != nil {
						log.Fatalf("db: %v\n", err)
					}
				}
			}

			// create the database
			log.Printf("db: create: creating database in %s\n", argsDb.paths.database)
			err := sqlite.Create(argsDb.paths.database, context.Background())
			if err != nil {
				log.Fatalf("db: create: %v\n", err)
			}

			log.Printf("db: created %q\n", argsDb.paths.database)
		},
	}
)
