package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Arafo/bitrise-step-create-github-pull-request/github"
	githubrepos "github.com/Arafo/bitrise-step-create-github-pull-request/githubrepo"
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
	return fmt.Sprintf("bot-tokens-%s", randomString(4))
}

func main() {

	// var conf = Config{
	// 	AuthToken:              "",
	// 	RepositoryURL:          "https://github.com/adimobile/ios-native-adidas-design-language",
	// 	NewBranch:              newBranchName(),
	// 	TargetBranch:           "main",
	// 	SourceFiles:            "step.yml",
	// 	PullRequestTitle:       "[Automatic] Generated new token JSON files",
	// 	PullRequestDescription: "This is an automatic PR generated when tokens have changed in .com repository. Accept & Merge this PR to use the latest values.",
	// 	APIBaseURL:             "https://api.github.com/",
	// 	BotName:                "svc-selfsigning",
	// 	BotEmail:               "bot@mail.com",
	// 	CommitMessage:          "new token resources files",
	// }

	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	conf.NewBranch = newBranchName()
	owner, repo := ownerAndRepo(conf.RepositoryURL)

	print("Using ")
	stepconf.Print(conf)
	println("")

	githubClient := getClient(conf)
	pullRequestClient := githubrepos.NewGithubRepository(*githubClient, repo, owner)

	err := cleanup(pullRequestClient, conf)

	if err != nil {
		log.Errorf("Failed to do housekeeping: %s\n", err)
		os.Exit(1)
	}

	err = createNewPr(pullRequestClient, conf)

	if err != nil {
		log.Errorf("Failed to create pull request: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)

}

func getClient(conf Config) *github.GithubClient {
	if conf.APIBaseURL == "" {
		return github.NewClient(string(conf.AuthToken))
	} else {
		return github.NewEnterpriseClient(conf.APIBaseURL, string(conf.AuthToken))
	}
}

// Cleanup will fetch all opened PRs to our target branch from our BOT and close them
func cleanup(repo *githubrepos.GithubRepository, config Config) error {

	fmt.Println("Cleaning up existing PRs:")
	fmt.Println("-------------------------")

	prs, err := repo.FetchPullRequests(config.NewBranch, config.TargetBranch)

	if err != nil {
		log.Errorf("Failed to fetch existing PRs: %s\n", err)
		return err
	}

	fmt.Println("Found", len(prs), "opened PR(s)")

	prsIdsWeAreIntrestedIn := []int{}
	for i := range prs {
		pr := prs[i]
		// if *pr.GetUser().Login == config.BotName {
		prsIdsWeAreIntrestedIn = append(prsIdsWeAreIntrestedIn, *pr.Number)
		// }
	}

	fmt.Println("Of which", len(prsIdsWeAreIntrestedIn), "need to be closed")
	fmt.Println("Closing PRs with Ids:", prsIdsWeAreIntrestedIn)

	for i := range prsIdsWeAreIntrestedIn {
		prId := prsIdsWeAreIntrestedIn[i]
		fmt.Print("Closing: ", prId, "...")
		err := repo.ClosePullRequest(prId)
		if err != nil {
			fmt.Print("Failed!\n")
			log.Errorf("Failed because of: %s\n", err)
		} else {
			fmt.Print("Ok!\n")
		}
	}

	fmt.Println("Done with cleanup")
	fmt.Println("")

	return nil

}

func createNewPr(repo *githubrepos.GithubRepository, config Config) error {

	fmt.Println("Creating new PR:")
	fmt.Println("----------------")

	targetBranch := config.TargetBranch
	newBranch := config.NewBranch
	sourceFiles := config.SourceFiles
	pullRequestTitle := config.PullRequestTitle
	pullRequestDescription := config.PullRequestDescription

	fmt.Printf("Getting last commit on target (%s)\n", targetBranch)

	sha, err := repo.GetLastCommit(config.TargetBranch)
	if err != nil {
		log.Errorf("Unable to get last commit SHA: %s\n", err)
		return err
	}
	fmt.Println("Got last commit", sha)

	fmt.Printf("Creating new branch (%s) from base (%s) at SHA (%s)\n", newBranch, targetBranch, sha)
	err = repo.CreateBranch(newBranch, targetBranch, sha)
	if err != nil {
		log.Errorf("Unable to create new branch: %s\n", err)
		return err
	}

	fmt.Println("New branch created", newBranch)

	fmt.Println("Commiting changs to branch...")
	err = repo.CreateCommit(targetBranch, newBranch, sourceFiles, config.BotName, config.BotEmail, config.CommitMessage)
	if err != nil {
		log.Errorf("Unable to create new branch: %s\n", err)
		return err
	}

	fmt.Println("Creating pull request...")
	number, err := repo.CreatePullRequest(pullRequestTitle, pullRequestDescription, newBranch, targetBranch)
	if err != nil {
		log.Errorf("Github API call failed when creating a Pull Request: %w\n", err)
		return err
	} else {
		fmt.Println("New pull request created id:", number)
	}

	return nil

}
