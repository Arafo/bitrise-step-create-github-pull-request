package main

import (
	"fmt"

	"github.com/Arafo/bitrise-step-create-github-pull-request/github"
	"github.com/bitrise-io/go-utils/log"
	gogithub "github.com/google/go-github/v42/github"
)

// Fetches all open prs to the target branch and returns the pr ids
func FetchAllOpenPrIds(githubClient github.GithubClient, conf Config) ([]int, error) {

	owner, repo := ownerAndRepo(conf.RepositoryURL)

	var filter = gogithub.PullRequestListOptions{
		Base:  conf.TargetBranch,
		State: "open",
	}

	fmt.Println("Looking for prs...")
	prs, _, err := githubClient.PullRequests.List(githubClient.Context, owner, repo, &filter)

	if err != nil {
		log.Errorf("Unable to list existing prs: %s\n", err)
		return nil, err
	}

	fmt.Println("Found", len(prs), "opened PR(s)")

	prsIdsWeAreIntrestedIn := []int{}
	for i := range prs {
		pr := prs[i]
		if *pr.GetUser().Login == conf.BotName {
			prsIdsWeAreIntrestedIn = append(prsIdsWeAreIntrestedIn, *pr.Number)
		}
	}

	fmt.Println("Of which", len(prsIdsWeAreIntrestedIn), "need(s) to be closed")
	fmt.Println("PR ids to close:", prsIdsWeAreIntrestedIn)

	return prsIdsWeAreIntrestedIn, nil

}

// Closes a list of PR Ids
func CloseOpenedPrs(githubClient github.GithubClient, conf Config, ids []int) {

	owner, repo := ownerAndRepo(conf.RepositoryURL)

	for _, prId := range ids {

		fmt.Print("Closing pull request with id: ", prId, "... ")

		err := closePr(githubClient, owner, repo, prId)

		if err != nil {
			fmt.Print("Failed\n")
			log.Errorf("Failed to close PR with id: %d error: %s\n", prId, err)
		} else {
			fmt.Print("Success\n")
		}

	}

}

func closePr(githubClient github.GithubClient, owner string, repo string, id int) error {

	targetState := "closed"

	var edit = gogithub.PullRequest{
		State: &targetState,
	}

	_, _, err := githubClient.PullRequests.Edit(githubClient.Context, owner, repo, id, &edit)

	return err

}
