package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

func main() {
	cfg, err := Env()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	list, _, err := client.Issues.ListByRepo(ctx, cfg.Owner, cfg.Repo, &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range list {
		if e.Body != nil {
			cmdstring := ParseIssue(*e.Body)
			if cmdstring == ""{
				continue
			}
			cmd := exec.Command("bash", "-c", cmdstring)
			output, err := cmd.Output()
			if err == nil {
				_, _, err = client.Issues.Edit(ctx, cfg.Owner, cfg.Repo, *e.Number, &github.IssueRequest{
					State: prt("closed"),
				})
				if err != nil {
					log.Println(err)
				}
				_, _, err = client.Issues.CreateComment(ctx, cfg.Owner, cfg.Repo, *e.Number, &github.IssueComment{
					Body: prt(MustWithTemplate(completed, nil,
						"command", func() interface{} { return cmdstring },
						"output", func() interface{} { return string(output) },
					)),
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func prt(s string) *string {
	return &s
}
