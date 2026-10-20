package storages

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

const githubGistBaseURL = "https://api.github.com/gists/"

type GithubGistStorage struct{}

func init() { Stoarges[GithubGistStorage{}.Name()] = GithubGistStorage{} }

func (s GithubGistStorage) Name() string { return "github_gist" }

func (s GithubGistStorage) Save(
	packages []packagemanagers.Package,
	config map[string]any,
) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	pkgs, err := json.Marshal(packages)
	if err != nil {
		return err
	}

	gistID := utils.GetFromConfig[string]("gist_id", config)
	filename := utils.GetFromConfig[string]("filename", config)
	auth := utils.GetFromConfig[string]("auth", config)

	bodyStruct := struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}{
		Files: map[string]struct {
			Content string `json:"content"`
		}{
			filename: {
				Content: string(pkgs),
			},
		},
	}

	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		githubGistBaseURL+gistID,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	request.Header.Set("Authorization", "Bearer "+auth)

	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// Ignoreing error because we cant do anyting anyways
		resBody, _ := io.ReadAll(res.Body)

		return errors.New("api returned " + res.Status + ": " + string(resBody))
	}

	return nil
}

func (s GithubGistStorage) Get(
	config map[string]any,
) ([]packagemanagers.Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	gistID := utils.GetFromConfig[string]("gist_id", config)
	filename := utils.GetFromConfig[string]("filename", config)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		githubGistBaseURL+gistID,
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// Ignoreing error because we cant do anyting anyways
		resBody, _ := io.ReadAll(res.Body)

		return nil, errors.New("api returned " + res.Status + ": " + string(resBody))
	}

	body := struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	file, ok := body.Files[filename]
	if !ok {
		return nil, errors.New("No such file with the given filename")
	}

	packages := []packagemanagers.Package{}

	if err := json.NewDecoder(
		strings.NewReader(file.Content),
	).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}
