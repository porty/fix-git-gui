package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/go-ini/ini"
)

func findGitConfig(basePath string) (string, error) {
	basePath = path.Clean(basePath)
	if basePath == "" || basePath == "/" {
		return "", os.ErrNotExist
	}

	configFile := path.Join(basePath, ".git", "config")

	stat, err := os.Stat(configFile)
	if err == nil {
		return configFile, nil
	}

	if os.IsNotExist(err) {
		return findGitConfig(path.Join(basePath, ".."))
	} else if err != nil {
		return "", errors.New("failed to stat file: " + err.Error())
	}

	if stat.IsDir() {
		return "", fmt.Errorf("unexpected directory found: %q", configFile)
	}

	return configFile, nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir = path.Clean(dir)
	gitConfig, err := findGitConfig(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Failed to find git directory")
			os.Exit(1)
		}
		panic(err)
	}

	cfg, err := ini.Load(gitConfig)
	if err != nil {
		panic(err)
	}

	if section, _ := cfg.GetSection("gui"); section == nil {
		return
	}
	cfg.DeleteSection("gui")

	// its been a while since I've seen globals :(
	ini.PrettyFormat = false
	ini.PrettyEqual = true
	ini.PrettySection = false
	cfg.SaveToIndent(gitConfig, "\t")
}
