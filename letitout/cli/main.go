package main

import (
	"fmt"
	"github.com/wolfulus/letitout/letitout/cli/cmd"
	"github.com/wolfulus/letitout/letitout/inlets"
	"os"
)

func main() {
	if inlets.Ok() == false {
		fmt.Println("Failed to execute inlets. Check whether it's installed or use 'lio inlets' to download it.")
		os.Exit(1)
	}

	err := cmd.Execute()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}