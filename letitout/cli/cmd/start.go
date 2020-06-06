package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
	. "github.com/wolfulus/letitout/letitout"
	"github.com/wolfulus/letitout/letitout/inlets"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts tunneling a project to an exit server.",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := "default"
		if len(args) > 0 {
			projectName = args[0]
		}

		p := GetProject(projectName)

		upstream, _ := cmd.Flags().GetString("upstream")
		if upstream == "" {
			upstream = p.Upstream
		}

		hostname, _ := cmd.Flags().GetString("hostname")
		if hostname == "" {
			hostname = p.Hostname
		}

		server, _ := cmd.Flags().GetString("server")
		if server == "" {
			server = p.Server
		}

		s := GetServer(server)

		UpdateDns(hostname, s)

		fmt.Printf("Starting inlets to tunnel %s to https://%s/\n", upstream, hostname)
		inlets.Tunnel(s.Address, s.Token, hostname, upstream)
	},
}

func init() {
	startCmd.PersistentFlags().String("server", "", "Overrides the server value.")
	startCmd.PersistentFlags().String("upstream", "", "Overrides the upstream value.")
	startCmd.PersistentFlags().String("hostname", "", "Overrides the hostname value.")
	rootCmd.AddCommand(startCmd)
}
