package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
)

func getCliArgs() (string, []string) {
	var apiToken, repoListString, githubApiBase string

	flag.StringVar(&apiToken, "t", "", "API token")
	flag.StringVar(&repoListString, "r", "", "Comma separated list of repos")
	flag.StringVar(&githubApiBase, "g", "https://api.github.com/", "URL for Github's base v3 api")
	flag.Parse()

	repoList := strings.Split(repoListString, ",")

	return apiToken, repoList
}

func validateCliArgs(apiToken string, repoList []string) bool {
	errors := false

	// todo: this does not work... repoList has a length of 1 even when no repos are provided... look into why this is
	if len(repoList) == 0 {
		errors = true
		fmt.Println("Please provide repos to parse")
	}

	if apiToken == "" {
		errors = true
		fmt.Println("Please provide a personal access token: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}

	return errors
}

func main() {
	apiToken, repoList, githubApiBase := getCliArgs()
	errors := validateCliArgs(apiToken, repoList)

	if errors {
		return
	}

	// todo: this is just for testing...
	return

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := ghClient.Repositories.List(ctx, "", nil)
	fmt.Println(repos)
	fmt.Println(err)
}
