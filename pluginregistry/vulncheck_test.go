package pluginregistry

import (
	"testing"

	"github.com/mheers/imagesumdb/config"
	"github.com/mheers/imagesumdb/image"
	"github.com/stretchr/testify/require"
)

func TestVulncheck(t *testing.T) {
	vulncheck := GetVulncheck()
	dir := t.TempDir()
	file := dir + "/testdb.yml"

	registry, repository, tag := "", "alpine", "3.10"

	img := image.NewImage(&config.Config{
		DBFile: file,
	}, registry, repository, tag)

	report, err := vulncheck.Scan(img)
	require.NoError(t, err)
	img.SetVulncheck(image.NewVulncheck(report))
}
