package kep

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v37/github"
	"github.com/salaxander/kepctl/pkg/auth"
	"golang.org/x/oauth2"
)

const owner = "kubernetes"
const repo = "enhancements"

var c *github.Client

func init() {
	client := oauth2.NewClient(context.Background(), &auth.TokenSource{})
	c = github.NewClient(client)
}

type KEP struct {
	IssueNumber string
	Milstone    string
	SIG         string
	Stage       string
	Title       string
	URL         string

	Tracked bool
}

func Get(kepNumber string) *KEP {
	issueInt, err := strconv.Atoi(kepNumber)
	if err != nil {

	}
	issue, _, err := c.Issues.Get(context.Background(), owner, repo, issueInt)
	if err != nil {

	}
	kep := issueToKEP(issue)
	if issue.Milestone != nil {
		kep.Milstone = *issue.Milestone.Title
	}

	return kep
}

func List(milestone string, sig string, stage string, tracked bool) []*KEP {
	var keps []*KEP
	var allIssues []*github.Issue
	opt := &github.IssueListByRepoOptions{}
	for {
		issues, resp, err := c.Issues.ListByRepo(context.Background(), owner, repo, opt)
		if err != nil {
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	if milestone != "" {
		allIssues = filterMilestone(milestone, allIssues)
	}
	if sig != "" {
		allIssues = filterSIG(sig, allIssues)
	}
	if stage != "" {
		allIssues = filterStage(stage, allIssues)
	}
	if tracked {
		allIssues = filterTracked(allIssues)
	}
	for i := range allIssues {
		keps = append(keps, issueToKEP(allIssues[i]))
	}
	return keps
}

func issueToKEP(issue *github.Issue) *KEP {
	numStr := strconv.Itoa(*issue.Number)
	kep := &KEP{
		IssueNumber: numStr,
		Title:       *issue.Title,
		URL:         *issue.HTMLURL,
	}
	if issue.Milestone != nil {
		kep.Milstone = *issue.Milestone.Title
	}
	for i := range issue.Labels {
		if strings.Contains(*issue.Labels[i].Name, "sig") {
			kep.SIG = *issue.Labels[i].Name
		}
		if strings.Contains(*issue.Labels[i].Name, "stage") {
			kep.Stage = *issue.Labels[i].Name
		}
		if *issue.Labels[i].Name == "tracked/yes" {
			kep.Tracked = true
		}
	}

	return kep
}

func filterMilestone(milestone string, issues []*github.Issue) []*github.Issue {
	result := []*github.Issue{}
	for i := range issues {
		if issues[i].Milestone != nil {
			if *issues[i].Milestone.Title == milestone {
				result = append(result, issues[i])
			}
		}
	}
	return result
}

func filterSIG(sig string, issues []*github.Issue) []*github.Issue {
	result := []*github.Issue{}
	for i := range issues {
		for l := range issues[i].Labels {
			if *issues[i].Labels[l].Name == fmt.Sprintf("sig/%s", sig) {
				result = append(result, issues[i])
			}
		}
	}
	return result
}

func filterStage(stage string, issues []*github.Issue) []*github.Issue {
	result := []*github.Issue{}
	for i := range issues {
		for l := range issues[i].Labels {
			if *issues[i].Labels[l].Name == fmt.Sprintf("stage/%s", stage) {
				result = append(result, issues[i])
			}
		}
	}
	return result
}

func filterTracked(issues []*github.Issue) []*github.Issue {
	result := []*github.Issue{}
	for i := range issues {
		for l := range issues[i].Labels {
			if *issues[i].Labels[l].Name == "tracked/yes" {
				result = append(result, issues[i])
			}
		}
	}
	return result
}
