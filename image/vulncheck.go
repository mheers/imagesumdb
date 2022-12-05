package image

import (
	"context"
	"time"

	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/mheers/imagesumdb/trivyhelper"
)

type Vulncheck struct {
	LastChecked time.Time `yaml:"lastchecked"`
	LastResult  string    `yaml:"lastresult"`
	Unknown     int       `yaml:"unknown"`
	Low         int       `yaml:"low"`
	Medium      int       `yaml:"medium"`
	High        int       `yaml:"high"`
	Critical    int       `yaml:"critical"`
}

func (v *Vulncheck) Safe() bool {
	return v.Critical == 0 && v.High == 0
}

func (v *Vulncheck) Total() int {
	return v.Unknown + v.Low + v.Medium + v.High + v.Critical
}

func NewVulncheck(report *types.Report) *Vulncheck {
	vc := &Vulncheck{
		LastChecked: time.Now(),
	}

	if len(report.Results[0].Vulnerabilities) > 0 {
		vc.LastResult = report.Results[0].Vulnerabilities[0].VulnerabilityID

		severityCount := countSeverities(report.Results[0].Vulnerabilities)
		vc.Critical = severityCount["CRITICAL"]
		vc.High = severityCount["HIGH"]
		vc.Medium = severityCount["MEDIUM"]
		vc.Low = severityCount["LOW"]
		vc.Unknown = severityCount["UNKNOWN"]
	}
	return vc
}

func (i *Image) Scan() (*types.Report, error) {
	err := trivyhelper.InitTrivyHelper()
	if err != nil {
		return nil, err
	}

	report, err := trivyhelper.ScanImage(context.Background(), i.RepoTag())
	if err != nil {
		return nil, err
	}

	i.vulncheck = NewVulncheck(report)
	return report, nil
}

func (i *Image) Safe() bool {
	if i.vulncheck == nil {
		return false
	}
	return i.vulncheck.Safe()
}

func countSeverities(vulns []types.DetectedVulnerability) map[string]int {
	severityCount := map[string]int{}
	for _, v := range vulns {
		severityCount[v.Severity]++
	}
	return severityCount
}
