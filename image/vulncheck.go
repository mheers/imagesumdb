package image

import (
	"time"

	"github.com/aquasecurity/trivy/pkg/types"
)

type Vulncheck struct {
	LastChecked time.Time `yaml:"lastchecked"`
	Results     []*Vuln   `yaml:"results"`
	Unknown     int       `yaml:"unknown"`
	Low         int       `yaml:"low"`
	Medium      int       `yaml:"medium"`
	High        int       `yaml:"high"`
	Critical    int       `yaml:"critical"`
}

type Vuln struct {
	CVE      string `yaml:"id"`
	Title    string `yaml:"title"`
	PkgName  string `yaml:"pkgname"`
	Severity string `yaml:"severity"`
}

func NewVuln(vuln types.DetectedVulnerability) (v *Vuln) {
	v = &Vuln{
		CVE:      vuln.VulnerabilityID,
		Title:    vuln.Title,
		PkgName:  vuln.PkgName,
		Severity: vuln.Severity,
	}
	return
}

func NewVulncheck(report *types.Report) (vc *Vulncheck) {
	vc = &Vulncheck{
		LastChecked: time.Now(),
	}

	if len(report.Results) == 0 {
		return
	}

	if len(report.Results[0].Vulnerabilities) == 0 {
		return
	}

	vc.Results = make([]*Vuln, 0)

	for _, result := range report.Results {
		for _, vuln := range result.Vulnerabilities {
			vc.Results = append(vc.Results, NewVuln(vuln))
		}
	}

	for _, result := range report.Results {
		severityCount := countSeverities(result.Vulnerabilities)
		vc.Critical += severityCount["CRITICAL"]
		vc.High += severityCount["HIGH"]
		vc.Medium += severityCount["MEDIUM"]
		vc.Low += severityCount["LOW"]
		vc.Unknown += severityCount["UNKNOWN"]
	}
	return
}

func (v *Vulncheck) Safe() bool {
	return v.Critical == 0 && v.High == 0
}

func (v *Vulncheck) Total() int {
	return v.Unknown + v.Low + v.Medium + v.High + v.Critical
}

func (i *Image) Safe() bool {
	if i.I_vulncheck == nil {
		return false
	}
	return i.I_vulncheck.Safe()
}

func countSeverities(vulns []types.DetectedVulnerability) map[string]int {
	severityCount := map[string]int{}
	for _, v := range vulns {
		severityCount[v.Severity]++
	}
	return severityCount
}
