package main

import (
	"context"

	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/mheers/imagesumdb/image"
	"github.com/mheers/imagesumdb/plugin/imagesumdb-plugin-vulncheck/trivyhelper"
)

func Scan(i *image.Image) (*types.Report, error) {
	err := trivyhelper.InitTrivyHelper()
	if err != nil {
		return nil, err
	}

	report, err := trivyhelper.ScanImage(context.Background(), i.RegistryRepositoryTagPlain())
	if err != nil {
		return nil, err
	}

	return report, nil
}
