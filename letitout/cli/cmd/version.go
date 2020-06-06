package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wolfulus/letitout/letitout/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of LetItOut",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: ", version.Version)
	},
}