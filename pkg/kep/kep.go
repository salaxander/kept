package kep

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v37/github"
	"github.com/salaxander/kept/pkg/auth"
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
	Title             string            `yaml:"title"`
	IssueNumber       string            `yaml:"kep-number"`
	SIG               string            `yaml:"owning-sig"`
	ParticipatingSIGs []string          `yaml:"participating-sigs"`
	Status            string            `yaml:"status"`
	CreationDate      string            `yaml:"creation-date"`
	Stage             string            `yaml:"stage"`
	LatestMilestone   string            `yaml:"latest-milestone"`
	Milestone         map[string]string `yaml:"milestone"`
	URL               string

	Tracked bool
}

func Get(kepNumber string) *KEP {
	issueInt, _ := strconv.Atoi(kepNumber)
	issue, _, _ := c.Issues.Get(context.Background(), owner, repo, issueInt)
	kep := issueToKEP(issue)

	return kep
}

func List(milestone string, sig string, stage string, tracked bool) []*KEP {
	var keps []*KEP
	var allIssues []*github.Issue
	opt := &github.IssueListByRepoOptions{}
	for {
		issues, resp, _ := c.Issues.ListByRepo(context.Background(), owner, repo, opt)
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
		kep.LatestMilestone = *issue.Milestone.Title
	}
	for i := range issue.Labels {
		if strings.Contains(*issue.Labels[i].Name, "sig") {
			kep.SIG = splitLabel(*issue.Labels[i].Name)
		}
		if strings.Contains(*issue.Labels[i].Name, "stage") {
			kep.Stage = splitLabel(*issue.Labels[i].Name)
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

func FindKEPYaml(kepNumber string) (string, error) {
	issueInt, _ := strconv.Atoi(kepNumber)
	issue, _, _ := c.Issues.Get(context.Background(), owner, repo, issueInt)
	kep := issueToKEP(issue)
	path := fmt.Sprintf("/keps/sig-%s/", kep.SIG)
	_, kepDirContent, _, err := c.Repositories.GetContents(context.Background(), "kubernetes", "enhancements", path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", err
	}
	for i := range kepDirContent {
		dir := kepDirContent[i]
		if strings.Contains(*dir.Name, kep.IssueNumber) {
			kepPath := fmt.Sprintf("%s/kep.yaml", *dir.Path)
			kepContentEncoded, _, _, err := c.Repositories.GetContents(context.Background(), "kubernetes", "enhancements", kepPath, &github.RepositoryContentGetOptions{})
			if err != nil {
				return "", err
			}
			return kepContentEncoded.GetContent()
		}
	}
	return "", nil
}

func splitLabel(label string) string {
	s := strings.Split(label, "/")
	return s[1]
}
