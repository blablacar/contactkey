package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type versionControlSystem interface {
	retrieveSha1ForProject(branch string) (string, error)
	diff(deployedSha1 string, sha1ToDeploy string) (*Changes, error)
}

func RetrieveSha1ForProject(vcs versionControlSystem, branch string) (string, error) {
	return vcs.retrieveSha1ForProject(branch)
}

func Diff(vcs versionControlSystem, deployedSha1 string, sha1ToDeploy string) (*Changes, error) {
	return vcs.diff(deployedSha1, sha1ToDeploy)
}

type Stash struct {
	Repository string
	Project    string
	User       string
	Password   string
	Url        string
}

type Changes struct {
	Commits []Commits
}

type Commits struct {
	DisplayId      string
	AuthorFullName string
	AuthorSlug     string
	Title          string
}

type StashResponse struct {
	Values []struct {
		Id        string `json:"id"`
		DisplayId string `json:"displayId"`
		Message   string `json:"message"`
		Author    struct {
			DisplayName string `json:"displayName"`
			Slug        string `json:"slug"`
		} `json:"author"`
		Parents []struct {
			Id string `json:"id"`
		} `json:"parents"`
	} `json:"values"`
}

func (s Stash) retrieveSha1ForProject(branch string) (string, error) {
	params := url.Values{}
	params.Add("until", branch)
	params.Add("limit", "1")
	stashResponse, err := s.getStashResponse(params)
	if err != nil {
		return "", err
	}

	if len(stashResponse.Values) == 0 || stashResponse.Values[0].Id == "" {
		return "", errors.New("Stash: Sha1 not found in the response")
	}

	return stashResponse.Values[0].Id, nil
}

func (s Stash) diff(deployedSha1 string, sha1ToDeploy string) (*Changes, error) {
	params := url.Values{}
	params.Add("since", deployedSha1)
	params.Add("until", sha1ToDeploy)
	stashResponse, err := s.getStashResponse(params)
	if err != nil {
		return nil, err
	}

	changes := new(Changes)
	for cnt := 0; cnt < len(stashResponse.Values); cnt++ {
		// We are sanitizing Stash
		if !strings.HasPrefix(stashResponse.Values[cnt].Message, "Merge") {
			continue
		}
		commits := Commits{}
		commits.Title = stashResponse.Values[cnt].Message
		commits.DisplayId = stashResponse.Values[cnt].DisplayId
		commits.AuthorFullName = stashResponse.Values[cnt].Author.DisplayName
		commits.AuthorSlug = stashResponse.Values[cnt].Author.Slug
		changes.Commits = append(changes.Commits, commits)
	}

	return changes, nil
}

func (s Stash) getStashResponse(params url.Values) (*StashResponse, error) {
	client := &http.Client{}
	baseUrl := fmt.Sprintf(
		"%s/rest/api/latest/projects/%s/repos/%s/commits?",
		s.Url,
		s.Repository,
		s.Project,
	)

	request, err := http.NewRequest("GET", baseUrl+params.Encode(), nil)
	request.SetBasicAuth(s.User, s.Password)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	stashResponse := new(StashResponse)
	err = json.Unmarshal(body, &stashResponse)
	if err != nil {
		return nil, err
	}

	return stashResponse, nil
}
