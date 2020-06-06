package letitout

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	CloudFlare map[string]CloudFlare `yaml:"cloudflare"`
	Servers    map[string]Server     `yaml:"servers"`
	Projects   map[string]Project    `yaml:"projects"`
}

var (
	config Config
)

func Initialize() {
	err := viper.UnmarshalKey("cloudflare", &config.CloudFlare)
	if err != nil {
		fmt.Println("Failed to load cloudflare configs:", err)
		os.Exit(1)
	}

	err = viper.UnmarshalKey("servers", &config.Servers)
	if err != nil {
		fmt.Println("Failed to load server configs:", err)
		os.Exit(1)
	}

	err = viper.UnmarshalKey("projects", &config.Projects)
	if err != nil {
		fmt.Println("Failed to load project configs:", err)
		os.Exit(1)
	}
}
