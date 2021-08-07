package main

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token     string `envconfig:"INPUT_TOKEN"`
	RepoOwner string `envconfig:"INPUT_REPOSITORY"`
	Repo      string
	Owner     string
}

func Env() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	repoOwner := strings.Split(cfg.RepoOwner, "/")
	if len(repoOwner) != 2 {
		return Config{}, fmt.Errorf("INPUT_GITHUB_REPOSITORY not set correctly")
	}
	cfg.Repo, cfg.Owner = repoOwner[0], repoOwner[1]
	return cfg, nil
}
