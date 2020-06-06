package letitout

import (
	"fmt"
	"os"
)

type Server struct {
	Address string `yaml:"address"`
	Token string `yaml:"token"`
}

func GetServer(name string) *Server {
	server, ok := config.Servers[name]
	if ok == false {
		fmt.Printf("Server %s doesn't exist.", name)
		os.Exit(1)
	}

	return &server
}
