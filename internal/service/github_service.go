package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type GithubService struct {
	ApiUrl string
}

func NewGithubService() *GithubService {
	return &GithubService{
		ApiUrl: "https://api.github.com/users",
	}
}

type gitHubSearchUserResponse struct {
	AvatarURL string `json:"avatar_url"`
}

func (g *GithubService) FetchUserGitHubAvatar(profileURL string) (string, error) {
	githubURL, err := url.Parse(profileURL)
	if err != nil {
		return "", err
	}
	username := path.Base(githubURL.Path)
	resp, err := http.Get(fmt.Sprintf("%s/%s", g.ApiUrl, username))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var u gitHubSearchUserResponse
	err = json.Unmarshal(b, &u)
	if err != nil {
		return "", err
	}
	return u.AvatarURL, nil
}
