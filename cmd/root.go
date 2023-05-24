package main

import (
	"github.com/mheers/imagesumdb/helpers"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "imagesumdb",
		Short: "imagesumdb is a tool to manage images in a database",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			helpers.PrintInfo()
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func main() {
	err := Execute()
	if err != nil {
		panic(err)
	}
}
