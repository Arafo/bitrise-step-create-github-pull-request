package github

import (
	"context"
	"net/http"
	"os"

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
	tc := newAuthTokenClient(accessToken)
	return &GithubClient{ctx, gogithub.NewClient(tc)}
}

func NewEnterpriseClient(baseURL string, accessToken string) *GithubClient {
	ctx := context.Background()
	tc := newAuthTokenClient(accessToken)
	enterpriseClient, err := gogithub.NewEnterpriseClient(baseURL, baseURL, tc)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	return &GithubClient{ctx, enterpriseClient}
}

func newAuthTokenClient(accessToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}
