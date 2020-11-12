package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/contractshark/vypie/api"
	"github.com/contractshark/vypie/cli"
	"github.com/contractshark/vypie/golang"
	"github.com/contractshark/vypie/javascript"
	"github.com/contractshark/vypie/language"
	"github.com/fatih/color"
)

func main() {
	// enable colored output on ci
	if os.Getenv("GITHUB_ACTIONS") != "" {
		color.NoColor = false
	}

	// check token
	_, ok := os.LookupEnv("SERIESCI_TOKEN")
	if !ok {
		panic(errors.New("cannot find SERIESCI_TOKEN environment variable"))
	}
	cli.Checkf("environment variable %s found\n", cli.Blue("SERIESCI_TOKEN"))

	// check programming language
	lang, err := language.Detect(".")
	if err != nil {
		panic(err)
	}

	// run programming language related checks
	switch lang {
	case language.Go:
		if err := golang.Run(); err != nil {
			panic(err)
		}
	case language.JavaScript:
		if err := javascript.Run(); err != nil {
			panic(err)
		}
	}

	repo, err := api.Repo()
	if err != nil {
		panic(err)
	}

	cli.Checkln("I'm done. See", cli.Blue(fmt.Sprintf("https://contractshark.com/%s", repo)))
}
