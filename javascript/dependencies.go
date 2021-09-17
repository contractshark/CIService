package javascript

import (
	"errors"
	"fmt"

	"github.com/contractshark/byzn/api"
	"github.com/contractshark/byzn/cli"
)

func dependencies(packageJSON map[string]interface{}) error {
	dependencies, ok := packageJSON["dependencies"].(map[string]interface{})
	if !ok {
		return errors.New("dependencies not found")
	}

	cli.Checkf("%s dependencies found\n", cli.Blue(len(dependencies)))

	// create series
	if err := api.CreateByzn(api.ByznDependencies); err != nil {
		return err
	}

	if err := api.Post(fmt.Sprintf("%d", len(dependencies)), api.ByznDependencies); err != nil {
		return err
	}

	return nil
}
