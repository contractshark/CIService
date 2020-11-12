package javascript

import (
	"fmt"

	"github.com/contractshark/vypie/api"
	"github.com/contractshark/vypie/cli"
	"github.com/contractshark/vypie/size"
)

func bundlesize() error {
	s, err := size.Directory("build")
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%fK", s)
	cli.Checkf("total size of \"build\" directory is %s\n", cli.Blue(str))

	// create series
	if err := api.CreateSeries(api.SeriesBundleSize); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.SeriesBundleSize); err != nil {
		return err
	}

	return nil
}
