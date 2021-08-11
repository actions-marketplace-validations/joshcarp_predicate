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
	log.Println(cfg)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	collaborators, _, err := client.Repositories.ListCollaborators(ctx, cfg.Owner, cfg.Repo, &github.ListCollaboratorsOptions{ListOptions: github.ListOptions{}})
	if err != nil {
		log.Fatal(err)
	}
	var isCollaborator bool
	for _, collaborator := range collaborators {
		if collaborator.GetLogin() == cfg.Actor {
			isCollaborator = true
			break
		}
	}
	if !isCollaborator {
		return
	}

	list, _, err := client.Issues.ListByRepo(ctx, cfg.Owner, cfg.Repo, &github.IssueListByRepoOptions{State: "open"})
	if err != nil {
		log.Fatal(err)
	}

	prComment := "\nFixes: "
	var payload Payload
	if cfg.Event == "pull-request" && cfg.EventPath != "" {
		payload, err = GetPayload(cfg.EventPath)
		log.Println(payload)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, e := range list {
		if e.Body != nil {
			cmdstring := ParseIssue(*e.Body)
			if cmdstring == "" {
				continue
			}
			cmd := exec.Command("bash", "-c", cmdstring)
			output, err := cmd.Output()
			if err == nil {
				if cfg.Event == "pull-request" {
					prComment += " " + e.GetURL()
					continue
				}
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
	payload.PullRequest.Number = 7
	if cfg.Event == "pull-request" {
		pr, _, err := client.PullRequests.Get(ctx, cfg.Owner, cfg.Repo, payload.PullRequest.Number)
		description := pr.GetBody()
		description += prComment
		pr.Body = &description
		if err != nil {
			log.Fatal(err)
		}
		_, _, err = client.PullRequests.Edit(ctx, cfg.Owner, cfg.Repo, payload.PullRequest.Number, pr)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func prt(s string) *string {
	return &s
}
