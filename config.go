package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token     string `envconfig:"INPUT_TOKEN"`
	RepoOwner string `envconfig:"INPUT_REPOSITORY"`
	Event     string `envconfig:"GITHUB_EVENT_NAME"`
	Actor     string `envconfig:"GITHUB_ACTOR"`
	Repo      string
	Owner     string
	EventPath string `envconfig:"GITHUB_EVENT_PATH"`
}

type Payload struct {
	PullRequest PullRequest `json:"payload"`
}
type PullRequest struct {
	Action string `json:"action"`
	Number int    `json:"number"`
}

func GetPayload(eventPath string) (Payload, error) {
	bytes, err := os.ReadFile(eventPath)
	if err != nil {
		return Payload{}, err
	}
	log.Println(string(bytes))
	var payload Payload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return Payload{}, err
	}
	return payload, nil
}

func Env() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	repoOwner := strings.Split(cfg.RepoOwner, "/")
	if len(repoOwner) != 2 {
		return Config{}, fmt.Errorf("INPUT_REPOSITORY not set correctly")
	}
	cfg.Owner, cfg.Repo = repoOwner[0], repoOwner[1]
	return cfg, nil
}
