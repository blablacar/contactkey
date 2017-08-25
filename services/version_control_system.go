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
	log "github.com/sirupsen/logrus"
)

type VersionControlSystem interface {
	RetrieveSha1ForProject(branch string) (string, error)
	Diff(deployedSha1 string, sha1ToDeploy string) (*Changes, error)
}

type Stash struct {
	Repository  string
	Project     string
	User        string
	Password    string
	Url         string
	Branch      string
	sha1MaxSize int
	Log         *log.Logger
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

func NewStash(cfg utils.StashConfig, manifest utils.StashManifest, logger *log.Logger) (*Stash, error) {
	if cfg.User == "" {
		return nil, errors.New("You need to define a user for stash in the config.")
	}

	if cfg.Password == "" {
		return nil, errors.New("You need to define an password for stash in the config.")
	}

	if cfg.Url == "" {
		return nil, errors.New("You need to define an url for stash in the config.")
	}

	if manifest.Repository == "" {
		return nil, errors.New("You need to define a repository for stash in the manifest.")
	}

	if manifest.Project == "" {
		return nil, errors.New("You need to define a project for stash in the manifest.")
	}

	if manifest.Branch == "" {
		return nil, errors.New("You need to define a branch for stash in the manifest.")
	}

	return &Stash{
		Repository:  manifest.Repository,
		Project:     manifest.Project,
		User:        cfg.User,
		Password:    cfg.Password,
		Url:         cfg.Url,
		Branch:      manifest.Branch,
		sha1MaxSize: cfg.Sha1MaxSize,
		Log:         logger,
	}, nil
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

	if s.sha1MaxSize > 0 {
		return stashResponse.Values[0].Id[0:s.sha1MaxSize], nil
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

	s.Log.WithFields(log.Fields{
		"changes": changes,
	}).Debug("Struct Changes")

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

	s.Log.WithFields(log.Fields{
		"fullPath":   baseUrl,
		"stashUrl":   s.Url,
		"project":    s.Project,
		"repository": s.Repository,
	}).Debug("Creating stash url")

	request, err := http.NewRequest("GET", baseUrl+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(s.User, s.Password)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	s.Log.WithFields(log.Fields{
		"statusCode": response.StatusCode,
	}).Debug("Stash response status code")

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Stash status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	s.Log.WithFields(log.Fields{
		"body": string(body),
	}).Debug("Response body from Stash")

	stashResponse := new(StashResponse)
	err = json.Unmarshal(body, &stashResponse)
	if err != nil {
		return nil, err
	}

	s.Log.WithFields(log.Fields{
		"stashResponse": stashResponse,
	}).Debug("Struct StashResponse")

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
