package javascript

import (
	"errors"
	"fmt"

	"github.com/contractshark/CIService/api"
	"github.com/contractshark/CIService/cli"
)

func dependencies(packageJSON map[string]interface{}) error {
	dependencies, ok := packageJSON["dependencies"].(map[string]interface{})
	if !ok {
		return errors.New("dependencies not found")
	}

	cli.Checkf("%s dependencies found\n", cli.Blue(len(dependencies)))

	// create shark
	if err := api.CreateShark(api.SharkDependencies); err != nil {
		return err
	}

	if err := api.Post(fmt.Sprintf("%d", len(dependencies)), api.SharkDependencies); err != nil {
		return err
	}

	return nil
}
