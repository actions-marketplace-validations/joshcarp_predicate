package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("INPUT_GH_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	list, _, err := client.Issues.ListByRepo(ctx, os.Getenv("INPUT_GH_OWNER"), os.Getenv("INPUT_GH_REPO"), &github.IssueListByRepoOptions{
		State: "open",
	})
	if err != nil {
		panic(err)
	}
	for _, e := range list {
		if e.Body != nil {
			cmdstring := ParseIssue(*e.Body)
			cmd := exec.Command("bash", "-c", cmdstring)
			output, err := cmd.Output()
			if err == nil {
				client.Issues.Edit(ctx, os.Getenv("INPUT_GH_OWNER"), os.Getenv("INPUT_GH_REPO"), *e.Number, &github.IssueRequest{
					State: prt("closed"),
				})
				fmt.Println(client.Issues.CreateComment(ctx, os.Getenv("INPUT_GH_OWNER"), os.Getenv("INPUT_GH_REPO"), *e.Number, &github.IssueComment{
					Body: prt(MustWithTemplate(completed, nil,
						"command", func() interface{} { return cmdstring },
						"output", func() interface{} { return string(output) },
					)),
				}))
			}

		}

	}

}

const completed = "This issue was closed because the command successfully executed.\nCommand:\n\n```\n{{command}}\n```\n\nOutput:\n\n```\n{{output}}\n```"

func prt(s string) *string {
	return &s
}

func ParseIssue(description string) string {
	var codeblock string
	description = strings.ReplaceAll(description, "\r", "")
	a := markdown.Parse([]byte(description), nil)
	ast.WalkFunc(a, func(node ast.Node, entering bool) ast.WalkStatus {
		leaf, ok := node.(*ast.CodeBlock)
		if !ok {
			return ast.GoToNext
		}
		if string(leaf.Info) != "predicate" {
			return ast.GoToNext
		}
		codeblock = string(leaf.Literal)
		return ast.Terminate
	})
	return codeblock
}

// MustWithTemplate calls WithTemplate and panics if err is not nil.
func MustWithTemplate(tmplstr string, data interface{}, funcs ...interface{}) string {
	val, err := WithTemplate(tmplstr, data, funcs...)
	if err != nil {
		panic(err)
	}
	return val
}


func WithTemplate(tmplstr string, data interface{}, funcs ...interface{}) (string, error) {
	funcmap := sprig.FuncMap()
	err := extraFuncs(funcmap, funcs...)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("").
		Funcs(map[string]interface{}(funcmap)).
		Parse(tmplstr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), err
}

func extraFuncs(m map[string]interface{}, funcs ...interface{}) error {
	if len(funcs)%2 != 0 {
		return fmt.Errorf("extra funcs should be even with form ['funcname', func...]")
	}
	for i := 0; i < len(funcs)-1; i += 2 {
		key, ok := funcs[i].(string)
		if !ok {
			return fmt.Errorf("key of wrong type, key should be string type")
		}
		val := funcs[i+1]
		m[key] = val
	}
	return nil
}
