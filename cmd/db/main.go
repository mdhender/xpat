// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/semver"
	"github.com/spf13/cobra"
	"log"
)

var (
	version = semver.Version{Major: 0, Minor: 1, Patch: 0}
)

func main() {
	log.SetFlags(log.Lshortfile)

	cmdRoot.AddCommand(cmdDb)
	cmdDb.PersistentFlags().StringVar(&argsDb.paths.database, "path", argsDb.paths.database, "path to the database file")

	cmdDb.AddCommand(cmdDbCreate)
	cmdDbCreate.AddCommand(cmdDbCreateDatabase)
	cmdDbCreateDatabase.Flags().BoolVar(&argsDbCreateDatabase.force, "force", false, "force the creation if the database exists")
	cmdDbCreateDatabase.Flags().StringVar(&argsDb.paths.database, "path", argsDb.paths.database, "path to the database file")
	if err := cmdDbCreateDatabase.MarkFlagRequired("path"); err != nil {
		log.Fatalf("path: %v\n", err)
	}

	cmdRoot.AddCommand(cmdVersion)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

var (
	cmdRoot = &cobra.Command{
		Use:   "db",
		Short: "Database management commands",
		Long:  `Run tasks against the database.`,
	}
)
