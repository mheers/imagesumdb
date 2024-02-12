// from github.com/arman1371/vulnerability-scanner-api

package trivyhelper

import (
	"bytes"
	"context"
	_ "embed"

	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/spf13/viper"
)

var (
	version     string                = "dev"
	globalFlags *flag.GlobalFlagGroup = flag.NewGlobalFlagGroup()
	imageFlags  *flag.Flags
)

//go:embed trivy-default.yaml
var trivyDefault []byte

func InitTrivyHelper() error {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(trivyDefault))
	if err != nil {
		return err
	}

	// // Initialize trivy logger
	// trivyLog.Logger = log.Logger
	// dlog.SetLogger(log.Logger)
	// flog.SetLogger(log.Logger)

	imageFlags = createImageFlags()
	return nil
}

func ScanImage(ctx context.Context, image string) (*types.Report, error) {

	options, err := imageFlags.ToOptions([]string{image})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	r, err := artifact.NewRunner(ctx, options)
	if err != nil {
		return nil, err
	}
	defer r.Close(ctx)

	var report types.Report
	if report, err = r.ScanImage(ctx, options); err != nil {
		return nil, err
	}

	report, err = r.Filter(ctx, options, report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func createImageFlags() *flag.Flags {
	reportFlagGroup := flag.NewReportFlagGroup()
	reportFlagGroup.ReportFormat = nil

	return &flag.Flags{
		CacheFlagGroup:         flag.NewCacheFlagGroup(),
		DBFlagGroup:            flag.NewDBFlagGroup(),
		ImageFlagGroup:         flag.NewImageFlagGroup(), // container image specific
		LicenseFlagGroup:       flag.NewLicenseFlagGroup(),
		MisconfFlagGroup:       flag.NewMisconfFlagGroup(),
		RemoteFlagGroup:        flag.NewClientFlags(), // for client/server mode
		RegoFlagGroup:          flag.NewRegoFlagGroup(),
		ReportFlagGroup:        reportFlagGroup,
		ScanFlagGroup:          flag.NewScanFlagGroup(),
		SecretFlagGroup:        flag.NewSecretFlagGroup(),
		VulnerabilityFlagGroup: flag.NewVulnerabilityFlagGroup(),
	}
}
