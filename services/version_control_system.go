package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/remyLemeunier/contactkey/utils"
)

type VersionControlSystem interface {
	RetrieveSha1ForProject(branch string) (string, error)
	Diff(deployedSha1 string, sha1ToDeploy string) (*Changes, error)
}

type Stash struct {
	Repository string
	Project    string
	User       string
	Password   string
	Url        string
	Branch     string
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

func NewStash(cfg utils.StashConfig, manifest utils.StashManifest) *Stash {
	return &Stash{
		Repository: manifest.Repository,
		Project:    manifest.Project,
		User:       cfg.User,
		Password:   cfg.Password,
		Url:        cfg.Url,
		Branch:     manifest.Branch,
	}
}

func (s Stash) RetrieveSha1ForProject(branch string) (string, error) {
	if branch == "" {
		branch = s.Branch
	}
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

func (s Stash) Diff(deployedSha1 string, sha1ToDeploy string) (*Changes, error) {
	params := url.Values{}
	params.Add("since", deployedSha1)
	params.Add("until", sha1ToDeploy)
	stashResponse, err := s.getStashResponse(params)
	if err != nil {
		return nil, err
	}

	changes := new(Changes)
	for cnt := 0; cnt < len(stashResponse.Values); cnt++ {
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
		s.Project,
		s.Repository,
	)

	request, err := http.NewRequest("GET", baseUrl+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

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

func (s *Stash) fill(config map[string]string) error {
	mandatoryIndexes := [5]string{"repository", "project", "user", "password", "url"}

	indexesMissing := make([]string, 0)
	for _, mandatoryIndex := range mandatoryIndexes {
		if config[mandatoryIndex] == "" {
			indexesMissing = append(indexesMissing, mandatoryIndex)
		}
	}

	if len(indexesMissing) > 0 {
		return errors.New(
			fmt.Sprintf(
				"Configuration '%s' for stash must be defined",
				strings.Join(indexesMissing[:], ","),
			))
	}

	s.Repository = config["repository"]
	s.Project = config["project"]
	s.User = config["user"]
	s.Password = config["password"]
	s.Url = config["url"]

	return nil
}
