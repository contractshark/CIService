package javascript

import (
	"os"
	"os/exec"
	"time"

	"github.com/contractshark/CIService/api"
	"github.com/contractshark/CIService/cli"
)

func duration() error {
	start := time.Now()

	cmd := exec.Command("npm", "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	elapsed := time.Since(start)

	cli.Checkf("build took %s\n", cli.Blue(elapsed))

	// create shark
	if err := api.CreateShark(api.SharkTime); err != nil {
		return err
	}

	// post value
	if err := api.Post(elapsed.String(), api.SharkTime); err != nil {
		return err
	}

	return nil
}
