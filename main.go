package main

import (
	"os"
	"strings"
	"encoding/json"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/Arafo/bitrise-step-create-github-pull-request/github"
)

type Config struct {
	AuthToken        		stepconf.Secret 	`env:"personal_access_token,required"`
	RepositoryURL    		string          	`env:"repository_url,required"`
	BaseBranch		 		string			 	`env:"base_branch,required"`	
	CommitBranch	 		string			 	`env:"commit_branch"`
	SourceFiles	 	 		string			 	`env:"source_files,required"`
	PullRequestTitle 		string			 	`env:"pull_request_title,required"`
	PullRequestDescription	string  		 	`env:"pull_request_description"`
	APIBaseURL       		string         	 	`env:"api_base_url,required"`
	Debug		       		bool         	 	`env:"debug"`
}

func ownerAndRepo(url string) (string, string) {
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "git@")
	paths := strings.FieldsFunc(url, func(r rune) bool { return r == '/' || r == ':' })
	return paths[1], strings.TrimSuffix(paths[2], ".git")
}

func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)

	owner, repo := ownerAndRepo(conf.RepositoryURL)
	baseBranch := conf.BaseBranch
	commitBranch := conf.CommitBranch
	sourceFiles := conf.SourceFiles
	pullRequestTitle := conf.PullRequestTitle
	pullRequestDescription := conf.PullRequestDescription

	var githubClient *github.GithubClient

	if conf.APIBaseURL == "" {
		githubClient = github.NewClient(string(conf.AuthToken))
	} else {
		githubClient = github.NewEnterpriseClient(conf.APIBaseURL, string(conf.AuthToken))
	}

	ref, err := githubClient.GetCommitBranchReference(owner, repo, baseBranch, commitBranch)
	if err != nil {
		log.Errorf("Unable to get/create the commit reference: %s\n", err)
	}

	tree, err := githubClient.GenerateTree(owner, repo, sourceFiles, ref)
	if err != nil {
		log.Errorf("Unable to create the tree based on the provided files: %s\n", err)
	}

	commit, err := githubClient.PushCommit(owner, repo, "authorName", "authorEmail", "commitMessage", ref, tree)
	if err != nil {
		log.Errorf("Unable to create the commit: %s\n", err)
	}


	pullRequest, err := githubClient.CreatePullRequest(owner, repo, pullRequestTitle, commitBranch, baseBranch, pullRequestDescription)
	if err != nil {
		log.Errorf("Github API call failed when creating a Pull Request: %w\n", err)
		os.Exit(1)
	}

	if conf.Debug {
		log.Successf("- Commit:\n%v\n- Pull Request:\n%v\n", prettyPrint(commit), prettyPrint(pullRequest))
	}

	os.Exit(0)
}
