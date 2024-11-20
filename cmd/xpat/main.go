// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/semver"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var (
	version = semver.Version{Major: 0, Minor: 1, Patch: 0}
)

func main() {
	log.SetFlags(log.Lshortfile)

	cmdRoot.AddCommand(cmdVersion)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

var (
	cmdRoot = &cobra.Command{
		Use:   "xpat",
		Short: "application web server",
		Long:  `Run a simple to-do web server.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("xpat: %s\n", version.String())

			mux := http.NewServeMux()

			// Your application routes
			mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, World!"))
			})

			// Wrap the mux with the embedded file middleware
			handler := embeddedAssetsHandler(mux)

			log.Println("Server is running on http://localhost:8080")
			if err := http.ListenAndServe(":8080", handler); err != nil {
				log.Fatal(err)
			}
		},
	}
)
