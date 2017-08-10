package services

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var apiStub *httptest.Server

func setupTestCase(t *testing.T) func(t *testing.T) {
	apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/rest/api/latest/projects/project/repos/repository/commits?limit=1&until=branch":
			raw, err := ioutil.ReadFile("./data/sha1.json")
			if err != nil {
				t.Error(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(raw)
		case "/rest/api/latest/projects/project/repos/repository/commits?since=b04ad09883d1858081702b8e2d80eb348ead9849&until=b0d5ca3e586d48cc6d3ad35f0e03dfc891e62752":
			raw, err := ioutil.ReadFile("./data/diff.json")
			if err != nil {
				t.Error(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(raw)
		}

	}))

	return func(t *testing.T) {
		// teardown
	}
}

func TestStashRetrieveSha1ForProject(t *testing.T) {
	tearDown := setupTestCase(t)
	defer tearDown(t)

	stash := Stash{
		"repository",
		"project",
		"user",
		"password",
		apiStub.URL,
		"defaultBranch",
		7,
	}

	sha1, err := stash.RetrieveSha1ForProject("branch")
	if err != nil {
		t.Errorf("Error shouldn't have been raised instead got %s", err.Error())
	}

	if sha1 != "dbddae9" {
		t.Errorf("The sha1 defined in body is dbddae9not sha1 %s", sha1)
	}
}

func TestDiff(t *testing.T) {
	tearDown := setupTestCase(t)
	defer tearDown(t)

	stash := Stash{
		"repository",
		"project",
		"user",
		"password",
		apiStub.URL,
		"defaultBranch",
		7,
	}

	changes, err := stash.Diff("b04ad09883d1858081702b8e2d80eb348ead9849", "b0d5ca3e586d48cc6d3ad35f0e03dfc891e62752")
	if err != nil {
		t.Errorf("Error shouldn't have been raised instead got %s", err.Error())
	}

	if len(changes.Commits) != 4 {
		t.Errorf("commits length is 4 instead got %d", len(changes.Commits))
	}

	if changes.Commits[0].AuthorFullName != "Fullname user" {
		t.Errorf("Error got %s", changes.Commits[0].AuthorFullName)
	}

	if changes.Commits[0].DisplayId != "b0d5ca3e586" {
		t.Errorf("Error got %s", changes.Commits[0].DisplayId)
	}
	if changes.Commits[0].AuthorSlug != "slug" {
		t.Errorf("Error got %s", changes.Commits[0].AuthorSlug)
	}
	if changes.Commits[0].Title != "Merge pull request #138 in repository/project from branch 1to master\n\n* commit '40acf891cee4bb64fe16c213e97333d83cb5f682':\n  Some comments" {
		t.Errorf("Error got %s", changes.Commits[0].Title)
	}

	if changes.Commits[2].AuthorFullName != "Fullname user2" {
		t.Errorf("Error got %s", changes.Commits[1].AuthorFullName)
	}
	if changes.Commits[2].DisplayId != "1c48ab0e9c1" {
		t.Errorf("Error got %s", changes.Commits[1].DisplayId)
	}
	if changes.Commits[2].AuthorSlug != "slug2" {
		t.Errorf("Error got %s", changes.Commits[1].AuthorSlug)
	}

	if changes.Commits[2].Title != "Merge pull request #152 in repository/project from branch2 to master\n\n* commit 'd436e3c2b0385afc38bf6fb9b29567ef9b9f226b':\n  Some comments 2" {
		t.Errorf("Error got %s", changes.Commits[1].Title)
	}
}

func TestFill(t *testing.T) {
	stash := Stash{}
	config := make(map[string]string)
	config["repository"] = "repository"
	config["project"] = "project"
	config["user"] = "user"
	config["password"] = "password"
	config["url"] = "http://127.0.0.1:8080"

	err := stash.fill(config)
	if err != nil {
		t.Errorf("Error shouldn't have been raised instead got %s", err.Error())
	}

	// Delete config and check if we raise and error
	delete(config, "repository")
	delete(config, "user")

	err = stash.fill(config)
	if err == nil {
		t.Error("An error should have been raised because we have removed part of the config.")
	}

	if err.Error() != "Configuration 'repository,user' for stash must be defined" {
		t.Errorf("Wrong error got %s", err.Error())
	}
}
