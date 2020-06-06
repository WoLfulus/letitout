package inlets

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func Ok() bool {
	cmd := exec.Command("inlets", "version")
	return cmd.Run() == nil
}

func Download() {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "inlets", "inlets")
	if err != nil {
		fmt.Println("Failed to fetch latest inlets release:", err)
		os.Exit(1)
	}

	version := release.TagName

	fmt.Println("Latest is:", version)

	file := ""
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
		file = fmt.Sprintf("https://github.com/inlets/inlets/releases/download/%s/inlets.exe", *version)
	} else if runtime.GOOS == "darwin" {
		file = fmt.Sprintf("https://github.com/inlets/inlets/releases/download/%s/inlets-darwin", *version)
	} else if runtime.GOARCH == "arm64" {
		file = fmt.Sprintf("https://github.com/inlets/inlets/releases/download/%s/inlets-arm64", *version)
	} else {
		file = fmt.Sprintf("https://github.com/inlets/inlets/releases/download/%s/inlets", *version)
	}

	download, err := http.Get(file)
	if err != nil {
		fmt.Println("Failed to download inlets:", err)
		os.Exit(1)
	}
	defer download.Body.Close()

	output, err := os.Create(fmt.Sprintf("inlets%s", extension))
	if err != nil {
		fmt.Println("Failed to write inlets executable:", err)
		os.Exit(1)
	}
	defer output.Close()

	_, err = io.Copy(output, download.Body)
	if err != nil {
		fmt.Println("Failed to copy download stream:", err)
		os.Exit(1)
	}
}

func Tunnel(address string, token string, hostname string, upstream string) {
	cmd := exec.Command("inlets",
		"client",
		"--print-token=false",
		"--remote", address,
		"--token", token, "--upstream",
		fmt.Sprintf("%s=%s", hostname, upstream),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to start inlets:", err)
		os.Exit(1)
	}

	/*
	if err := cmd.Wait(); err != nil {
		fmt.Println("Inlets failed:", err)
		os.Exit(1)
	}
	*/
}
