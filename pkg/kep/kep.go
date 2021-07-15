package kep

import (
	"context"
	"strconv"

	"github.com/google/go-github/v37/github"
)

const owner = "kubernetes"
const repo = "enhancements"

var c *github.Client

func init() {
	c = github.NewClient(nil)
}

type KEP struct {
	IssueNumber string
	Title       string
	URL         string
}

func Get(kepNumber string) *KEP {
	issueInt, err := strconv.Atoi(kepNumber)
	if err != nil {

	}
	issue, _, err := c.Issues.Get(context.Background(), owner, repo, issueInt)
	if err != nil {

	}
	return &KEP{
		IssueNumber: strconv.Itoa(*issue.Number),
		Title:       *issue.Title,
		URL:         *issue.HTMLURL,
	}
}

func List() []*KEP {
	var keps []*KEP
	issues, _, err := c.Issues.ListByRepo(context.Background(), owner, repo, &github.IssueListByRepoOptions{ListOptions: github.ListOptions{PerPage: 500}})
	if err != nil {

	}
	for i := range issues {
		keps = append(keps, issueToKEP(issues[i]))
	}
	return keps
}

func issueToKEP(issue *github.Issue) *KEP {
	numStr := strconv.Itoa(*issue.Number)
	return &KEP{
		IssueNumber: numStr,
		Title:       *issue.Title,
		URL:         *issue.HTMLURL,
	}
}
