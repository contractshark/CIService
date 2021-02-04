package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/contractshark/shark/api"
	"github.com/contractshark/shark/cli"
	"github.com/contractshark/shark/golang"
	"github.com/contractshark/shark/javascript"
	"github.com/contractshark/shark/language"
	"github.com/fatih/color"
)

func main() {
	// enable colored output on ci
	if os.Getenv("GITHUB_ACTIONS") != "" {
		color.NoColor = false
	}

	// check token
	_, ok := os.LookupEnv("CONTRACT_SHARK_TOKEN")
	if !ok {
		panic(errors.New("cannot find CONTRACT_SHARK_TOKEN environment variable"))
	}
	cli.Checkf("environment variable %s found\n", cli.Blue("CONTRACT_SHARK_TOKEN"))

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
