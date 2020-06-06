package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wolfulus/letitout/letitout/inlets"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(inletsCmd)
}

var inletsCmd = &cobra.Command{
	Use:   "inlets",
	Short: "Download inlets executable",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading inlets...")
		inlets.Download()

		command := exec.Command("inlets", "version")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Println("Inlets execution failed:", err)
			os.Exit(1)
		}
	},
}