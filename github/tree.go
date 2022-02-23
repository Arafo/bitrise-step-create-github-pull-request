package github

import (
	"errors"
	"io/ioutil"
	"strings"

	gogithub "github.com/google/go-github/v42/github"
)

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
