package db

import (
	"testing"

	"github.com/mheers/imagesumdb/config"
	"github.com/stretchr/testify/require"
)

func TestVulncheck(t *testing.T) {
	dir := t.TempDir()
	file := dir + "/testdb.yml"
	db := NewDB(&config.Config{
		DBFile: file,
	})

	err := db.Add("test", "", "alpine", "3.10")
	require.NoError(t, err)

	err = db.Vulncheck()
	require.NoError(t, err)

	err = db.Write()
	require.NoError(t, err)
}
