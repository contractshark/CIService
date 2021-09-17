package javascript

import (
	"fmt"

	"github.com/contractshark/CIService/api"
	"github.com/contractshark/CIService/cli"
	"github.com/contractshark/CIService/size"
)

func bundlesize() error {
	s, err := size.Directory("build")
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%fK", s)
	cli.Checkf("total size of \"build\" directory is %s\n", cli.Blue(str))

	// create shark
	if err := api.CreateShark(api.SharkBundleSize); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.SharkBundleSize); err != nil {
		return err
	}

	return nil
}
