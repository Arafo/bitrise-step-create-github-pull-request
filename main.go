package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Arafo/bitrise-step-create-github-pull-request/github"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type Config struct {
	AuthToken              stepconf.Secret `env:"personal_access_token,required"`
	RepositoryURL          string          `env:"repository_url,required"`
	TargetBranch           string          `env:"target_branch,required"`
	NewBranch              string          `env:"commit_branch"`
	SourceFiles            string          `env:"source_files,required"`
	PullRequestTitle       string          `env:"pull_request_title,required"`
	PullRequestDescription string          `env:"pull_request_description"`
	BotName                string          `env:"bot_name,required"`
	BotEmail               string          `env:"bot_email,required"`
	CommitMessage          string          `env:"commit_message,required"`
	APIBaseURL             string          `env:"api_base_url,required"`
}

func ownerAndRepo(url string) (string, string) {
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "git@")
	paths := strings.FieldsFunc(url, func(r rune) bool { return r == '/' || r == ':' })
	return paths[1], strings.TrimSuffix(paths[2], ".git")
}

// func prettyPrint(i interface{}) string {
// 	s, _ := json.MarshalIndent(i, "", "\t")
// 	return string(s)
// }

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func newBranchName() string {
	return fmt.Sprintf("new-tokens-%s", randomString(6))
}

func main() {

	var conf = Config{
		AuthToken:              "",
		RepositoryURL:          "https://github.com/adimobile/ios-native-adidas-design-language",
		NewBranch:              newBranchName(),
		TargetBranch:           "main",
		SourceFiles:            "step.yml",
		PullRequestTitle:       "[Automatic] Generated new token JSON files",
		PullRequestDescription: "This is an automatic PR generated when tokens have changed in .com repository. Accept & Merge this PR to use the latest values.",
		APIBaseURL:             "https://api.github.com/",
		BotName:                "svc-selfsigning",
		BotEmail:               "bot@mail.com",
		CommitMessage:          "new token resources files",
	}

	// if err := stepconf.Parse(&conf); err != nil {
	// 	log.Errorf("Error: %s\n", err)
	// 	os.Exit(1)
	// }

	print("Using ")
	stepconf.Print(conf)
	println("")

	var githubClient *github.GithubClient

	if conf.APIBaseURL == "" {
		githubClient = github.NewClient(string(conf.AuthToken))
	} else {
		githubClient = github.NewEnterpriseClient(conf.APIBaseURL, string(conf.AuthToken))
	}

	fmt.Println("Cleaning up existing Pull Requests:")
	fmt.Println("----------------------------------")
	openedPrIds, err := FetchAllOpenPrIds(*githubClient, conf)
	if err != nil {
		log.Errorf("Failed to fetch all opened PRs error: %s\n", err)
		os.Exit(1)
	}
	CloseOpenedPrs(*githubClient, conf, openedPrIds)

	fmt.Println("")

	fmt.Println("Creating new Pull Request:")
	fmt.Println("--------------------------")
	err = CreateNewPr(*githubClient, conf)
	if err != nil {
		log.Errorf("Failed to create new pull request: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)

}
