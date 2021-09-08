package kep

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v37/github"

	"golang.org/x/oauth2"

	"github.com/salaxander/kept/pkg/auth"
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

func Get(kepNumber string) (*KEP, error) {
	var kep KEP

	// Set the KEP URL
	kep.URL = fmt.Sprintf("https://github.com/kubernetes/enhancements/issues/%s", kepNumber)

	// Get the KEP issue to determine SIG from label.
	issueInt, _ := strconv.Atoi(kepNumber)
	issue, _, err := c.Issues.Get(context.Background(), owner, repo, issueInt)
	if err != nil {
		return nil, err
	}

	// Determine the KEP's SIG.
	sig := findSIG(issue)

	// Get the KEP's raw kep.yaml content
	kepYaml, err := getKEPYaml(sig, kepNumber)
	if err != nil {
		return nil, err
	}
	// Unmarshall the yaml into the KEP object
	if err := yaml.Unmarshal([]byte(kepYaml), &kep); err != nil {
		return nil, err
	}

	return &kep, nil
}

func findSIG(issue *github.Issue) string {
	for i := range issue.Labels {
		label := issue.Labels[i]
		if strings.Contains(*issue.Labels[i].Name, "sig") {
			s := strings.Split(*label.Name, "/")
			return s[1]
		}
	}
	return ""
}

func getKEPYaml(sig, kepNumber string) (string, error) {
	path := fmt.Sprintf("/keps/sig-%s/", sig)
	_, kepDirContent, _, err := c.Repositories.GetContents(context.Background(), "kubernetes", "enhancements", path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", err
	}
	for i := range kepDirContent {
		dir := kepDirContent[i]
		if strings.Contains(*dir.Name, kepNumber) {
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
