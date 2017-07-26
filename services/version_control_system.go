package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"net/url"

	"github.com/labstack/gommon/log"
)

type versionControlSystem interface {
	retrieveSha1ForProject(branch string) string
	diff(deployedSha1 string, sha1ToDeploy string) Changes
}

func RetrieveSha1ForProject(vcs versionControlSystem, branch string) string {
	return vcs.retrieveSha1ForProject(branch)
}

func Diff(vcs versionControlSystem, deployedSha1 string, sha1ToDeploy string) Changes {
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

func (s Stash) retrieveSha1ForProject(branch string) string {

	params := url.Values{}
	params.Add("until", branch)
	params.Add("limit", "1")
	stashResponse := s.getStashResponse(params)

	if len(stashResponse.Values) == 0 || stashResponse.Values[0].Id == "" {
		log.Fatal("Stash: Sha1 not found in the response")
	}

	return stashResponse.Values[0].Id
}

func (s Stash) diff(deployedSha1 string, sha1ToDeploy string) Changes {
	params := url.Values{}
	params.Add("since", deployedSha1)
	params.Add("until", sha1ToDeploy)
	stashResponse := s.getStashResponse(params)

	changes := Changes{}
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

	return changes
}

func (s Stash) getStashResponse(params url.Values) StashResponse {
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
		log.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var stashResponse StashResponse
	err = json.Unmarshal(body, &stashResponse)
	if err != nil {
		log.Fatal(err.Error())
	}

	return stashResponse
}
