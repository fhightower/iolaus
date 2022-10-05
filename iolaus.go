package main

import (
    "context"
	"fmt"
	"github.com/google/go-github/github"
    "golang.org/x/oauth2"
)

func main() {
	fmt.Println("hi!")
    ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
    fmt.Println(repos)
    fmt.Println(err)
}
