package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/contractshark/shark/cli"
)

// shark names
const (
	ContractCoverage      = "coverage"
	ContractFileSize      = "size"
	ContractTime          = "time"
	ContractBundleSize    = "bundlesize"
	ContractDependencies  = "dependencies"
	ContractPerformance   = "performance"
	ContractLint          = "lint conformnace"
	ContractPractices     = "practices"
	ContractKPI           = "kpi"
)

// Descriptions returns the description for a given shark name.
var Descriptions = map[string]string{
	ContractCoverage:      "Code coverage",
	ContractFileSize:      "File size",
	ContractTime:          "Build time",
	ContractBundleSize:    "Bundle size",
	ContractDependencies:  "Number of dependencies",
	ContractPerformance:   "CShark performance",
	ContractLint:          "CShark conformance ",
	ContractPractices:     "CShark best practices",
	ContractKPI:           "CShark KPI",
}

// get sha depending on ci environment.
func sha() (string, error) {
	// github actions
	sha, ok := os.LookupEnv("GITHUB_SHA")
	if ok {
		return sha, nil
	}

	// travis ci
	sha, ok = os.LookupEnv("TRAVIS_COMMIT")
	if ok {
		return sha, nil
	}

	// circle ci
	sha, ok = os.LookupEnv("CIRCLE_SHA1")
	if ok {
		return sha, nil
	}

	// default git
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Repo returns owner and repository in the form owner/repository.
func Repo() (string, error) {
	// github actions
	repo, ok := os.LookupEnv("GITHUB_REPOSITORY")
	if ok {
		return repo, nil
	}

	// travis ci
	repo, ok = os.LookupEnv("TRAVIS_REPO_SLUG")
	if ok {
		return repo, nil
	}

	// circle ci
	if _, ok = os.LookupEnv("CIRCLECI"); ok {
		username := os.Getenv("CIRCLE_PROJECT_USERNAME")
		reponame := os.Getenv("CIRCLE_PROJECT_REPONAME")
		return username + "/" + reponame, nil
	}

	return "", errors.New("cannot find repo in environment variables")
}

// Post something.
// currently only works with GitHub Actions.
func Post(value, shark string) error {
	// get commit hash
	s, err := sha()
	if err != nil {
		return err
	}

	data := url.Values{
		"value": {value},
		"sha":   {s},
	}

	// get repo in form owner/repo
	r, err := Repo()
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://contractshark.com/api/repos/%s/%s/combined", r, shark)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("SERIESCI_TOKEN")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cli.Checkf("post %s: status code: %s, body: %s\n", cli.Blue(shark), cli.Blue(res.StatusCode), cli.Blue(string(body)))

	return nil
}

// CreateContractRequest request
type CreateContractRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateContract creates a new shark.
func CreateContract(shark string) error {

	// create custom request
	data := CreateContractRequest{
		Name:        shark,
		Description: Descriptions[shark],
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(data); err != nil {
		return err
	}

	// get repo in form owner/repo
	r, err := Repo()
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://contractshark.com/api/repos/%s/shark", r)
	req, err := http.NewRequest(http.MethodPost, u, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("SERIESCI_TOKEN")))

	// send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		cli.Checkf("shark %s already exists\n", cli.Blue(shark))
	} else {
		cli.Checkf("shark %s created\n", cli.Blue(shark))
	}

	return nil
}
