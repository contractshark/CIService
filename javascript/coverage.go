package javascript

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/contractshark/CIService/api"
	"github.com/contractshark/CIService/cli"
	"github.com/contractshark/CIService/clover"
)

func coverage(packageJSON map[string]interface{}) error {
	// edit package.json and add clover coverage reporter
	packageJSON["jest"] = map[string][]string{
		"coverageReporters": []string{"clover"},
	}

	// override package.json temporarily
	b, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("package.json", b, 0644); err != nil {
		return err
	}

	// run code coverage
	args := []string{
		"test",
		"--",
		"--coverage",
		"--watchAll=false",
	}
	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// read coverage
	cloverXML, err := ioutil.ReadFile(filepath.Join("coverage", "clover.xml"))
	if err != nil {
		return err
	}

	var coverage cov.Coverage
	if err := xml.Unmarshal(cloverXML, &coverage); err != nil {
		return err
	}

	// covered statements
	cs := float64(coverage.Project.Metrics.CoveredStatements) / float64(coverage.Project.Metrics.Statements) * 100
	str := fmt.Sprintf("%.2f%%", cs)

	cli.Checkf("code coverage is %s\n", cli.Blue(str))

	// create shark
	if err := api.CreateShark(api.SharkCoverage); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.SharkCoverage); err != nil {
		return err
	}

	return nil

}
