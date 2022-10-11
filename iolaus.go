package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
)

type PR struct {
	owner  string
	repo   string
	number int
}

func getCliArgs() (string, []string) {
	var apiToken, prListString string

	flag.StringVar(&apiToken, "t", "", "API token")
	// flag.StringVar(&githubApiBase, "g", "https://api.github.com/", "URL for Github's base v3 api")
	flag.StringVar(&prListString, "prs", "", "Comma separated list of PRs")
	flag.Parse()

	prList := strings.Split(prListString, ",")

	return apiToken, prList
}

func validateCliArgs(apiToken string, prList []string) bool {
	errors := false

	// todo: this does not work... prList has a length of 1 even when no repos are provided... look into why this is
	if len(prList) == 0 {
		errors = true
		fmt.Println("Please provide a comma separated list of PRs to approve")
	}

	if apiToken == "" {
		errors = true
		fmt.Println("Please provide a personal access token: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}

	return errors
}

func processPRs(prList []string) []PR {
	var cleanedPRList []PR
	for _, v := range prList {
		// todo: process PRs here...
		owner := "fhightower"
		repo := "ioc-finder"
		number := 269
		cleanedPRList = append(cleanedPRList, PR{owner, repo, number})
	}
	return cleanedPRList
}

func main() {
	apiToken, prList := getCliArgs()
	errors := validateCliArgs(apiToken, prList)
	prs := processPRs(prList)

	if errors {
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	for _, pr := range prs {
		pr, _, err := ghClient.PullRequests.Get(ctx, pr.owner, pr.repo, pr.number)
		fmt.Println(pr)
		fmt.Println(err)
	}
	// todo: this is just for testing...
	return
}
