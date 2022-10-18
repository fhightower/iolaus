package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strconv"
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
	flag.StringVar(&prListString, "prs", "", "Comma separated list of PRs in the form: `{owner}/{repo}/pulls/{number}` (e.g. fhightower/ioc-finder/pull/271)")
	flag.Parse()

	prList := strings.Split(prListString, ",")

	return apiToken, prList
}

func validateCliArgs(apiToken string, prList []string) bool {
	errors := false

	if prList[0] == "" {
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
		eles := strings.Split(v, "/")
		owner := eles[0]
		repo := eles[1]
		number, _ := strconv.Atoi(eles[3])
		cleanedPRList = append(cleanedPRList, PR{owner, repo, number})
	}
	return cleanedPRList
}

func determineMergeableState(pr github.PullRequest) bool {
	// there are no clear docs for mergableState, but here the possible values are enumerated here:
	// https://github.com/octokit/octokit.net/issues/1763
	// and there are docs for the graphql equivalent (named 'mergestatestatus') here:
	// https://docs.github.com/en/graphql/reference/enums#mergestatestatus
	return pr.GetMergeableState() == "clean"
}

func main() {
	apiToken, prList := getCliArgs()
	errors := validateCliArgs(apiToken, prList)
	if errors {
		return
	}

	prs := processPRs(prList)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	for _, pr := range prs {
		thisPr, _, _ := ghClient.PullRequests.Get(ctx, pr.owner, pr.repo, pr.number)
		canMerge := determineMergeableState(*thisPr)
		if !canMerge {
			fmt.Printf("PR %v cannot be merged... please check this PR\n", pr)
			return
		}
		prBody := "🚀\n\n*(This PR was automatically approved and merged using [iolaus](https://github.com/fhightower/iolaus))*"
		prReview := github.PullRequestReviewRequest{Event: github.String("APPROVE"), Body: &prBody}
		_, _, errs := ghClient.PullRequests.CreateReview(ctx, pr.owner, pr.repo, pr.number, &prReview)
		if errs != nil {
			fmt.Println(errs)
			return
		}
		mergeOptions := github.PullRequestOptions{MergeMethod: "squash"}
		// the commit message is empty (""), so Github will use the default commit message
		_, _, err := ghClient.PullRequests.Merge(ctx, pr.owner, pr.repo, pr.number, "", &mergeOptions)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
}
