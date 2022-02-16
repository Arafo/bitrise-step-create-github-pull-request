package github

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"time"

	"github.com/bitrise-io/go-utils/log"
	gogithub "github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Context context.Context
	*gogithub.Client
}

func NewClient(accessToken string) *GithubClient {
	ctx := context.Background()
	tc := NewAuthTokenClient(accessToken)

	return &GithubClient{
		ctx,
		gogithub.NewClient(tc),
	}
}

func NewEnterpriseClient(baseURL string, accessToken string) *GithubClient {
	ctx := context.Background()
	tc := NewAuthTokenClient(accessToken)
	enterpriseClient, err := gogithub.NewEnterpriseClient(baseURL, baseURL, tc)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	return &GithubClient{
		ctx,
		enterpriseClient,
	}
}

func NewAuthTokenClient(accessToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}

func (c *GithubClient) GetCommitBranchReference(owner, repo string, baseBranch string, commitBranch string) (ref *gogithub.Reference, err error) {
	ref, _, err = c.Git.GetRef(c.Context, owner, repo, "refs/heads/" + commitBranch)
	if err == nil {
		return ref, nil
	}

	if commitBranch == baseBranch {
		return nil, errors.New("Commit branch does not exist but `base_branch` is the same as `commit_branch`")
	}

	var baseRef *gogithub.Reference
	if baseRef, _, err = c.Git.GetRef(c.Context, owner, repo, "refs/heads/" + baseBranch); err != nil {
		return nil, err
	}

	newRef := &gogithub.Reference{Ref: gogithub.String("refs/heads/" + commitBranch), Object: &gogithub.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = c.Git.CreateRef(c.Context, owner, repo, newRef)

	return ref, err
}

func (c *GithubClient) GenerateTree(owner, repo string, sourceFiles string, ref *gogithub.Reference) (tree *gogithub.Tree, err error) {
	entries := []*gogithub.TreeEntry{}

	for _, fileArg := range strings.Split(sourceFiles, "\n") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &gogithub.TreeEntry{Path: gogithub.String(file), Type: gogithub.String("blob"), Content: gogithub.String(string(content)), Mode: gogithub.String("100644")})
	}

	tree, _, err = c.Git.CreateTree(c.Context, owner, repo, *ref.Object.SHA, entries)
	return tree, err
}

func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

func (c *GithubClient) PushCommit(owner, repo string, authorName string, authorEmail string, commitMessage string, ref *gogithub.Reference, tree *gogithub.Tree) (commit *gogithub.Commit, err error) {
	parent, _, err := c.Repositories.GetCommit(c.Context, owner, repo, *ref.Object.SHA, nil)
	if err != nil {
		return nil, err
	}

	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := gogithub.CommitAuthor{Date: &date, Name: gogithub.String(authorName), Email: gogithub.String(authorEmail)}
	commitToPush := gogithub.Commit{Author: &author, Message: gogithub.String(commitMessage), Tree: tree, Parents: []*gogithub.Commit{parent.Commit}}
	newCommit, _, err := c.Git.CreateCommit(c.Context, owner, repo, &commitToPush)
	if err != nil {
		return nil, err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = c.Git.UpdateRef(c.Context, owner, repo, ref, false)
	return newCommit, err
}

func (c *GithubClient) CreatePullRequest(owner, repo string, title string, commitBranch string, baseBranch string, description string) (*gogithub.PullRequest, error) {
	pullRequestToCreate := gogithub.NewPullRequest{
							Title:               gogithub.String(title),
							Head:                gogithub.String(commitBranch),
							Base:                gogithub.String(baseBranch),
							Body:                gogithub.String(description),
							MaintainerCanModify: gogithub.Bool(true),
	}

	pullRequest, _, err := c.PullRequests.Create(c.Context, owner, repo, &pullRequestToCreate)

	return pullRequest, err
}