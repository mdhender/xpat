// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/xpat"
	"github.com/spf13/cobra"
)

var (
	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of this application",
		Long:  `Version of the server application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s (xpat %s)\n", version.String(), xpat.Version().String())
		},
	}
)
