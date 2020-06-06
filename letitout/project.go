package letitout

import (
	"fmt"
	"os"
)

type Project struct {
	Server   string `yaml:"server"`
	Hostname string `yaml:"host"`
	Upstream string `yaml:"upstream"`
}

func GetProject(name string) *Project {
	project, ok := config.Projects[name]
	if ok == false {
		fmt.Printf("Project %s doesn't exist.", name)
		os.Exit(1)
	}

	return &project
}
