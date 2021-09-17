package javascript

import (
	"fmt"

	"github.com/contractshark/byzn/api"
	"github.com/contractshark/byzn/cli"
	"github.com/contractshark/byzn/size"
)

func bundlesize() error {
	s, err := size.Directory("build")
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%fK", s)
	cli.Checkf("total size of \"build\" directory is %s\n", cli.Blue(str))

	// create series
	if err := api.CreateByzn(api.ByznBundleSize); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.ByznBundleSize); err != nil {
		return err
	}

	return nil
}
