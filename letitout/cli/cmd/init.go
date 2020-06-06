package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wolfulus/letitout/letitout"
	"os"
)

var (
	configFile string
	rootCmd = &cobra.Command{
		Use: "lio",
		Short: "LetItOut makes it easier to run inlets clients on development machines.",
	}
)

func Execute() error {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil || cmd == nil {
		args := append([]string{"start"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.letitout.yml)")
}

func abort(err interface{}) {
	fmt.Println("Error:", err)
	os.Exit(1)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			abort(err)
		}

		viper.SetConfigName(".letitout")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Failed to read config.")
		os.Exit(1)
	}

	letitout.Initialize()
}