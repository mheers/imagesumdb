package image

import (
	"testing"

	"github.com/mheers/imagesumdb/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifest(t *testing.T) {
	di := NewImage(nil, "docker.io/library", "alpine", "3.12")
	manifest, err := di.Manifest()
	require.NoError(t, err)
	require.NotEmpty(t, manifest)
}

func TestDigest(t *testing.T) {
	di := NewImage(nil, "docker.io/library", "alpine", "3.12")
	digest, err := di.Digest()
	require.NoError(t, err)
	require.NotEmpty(t, digest)
	assert.Equal(t, "sha256:cb64bbe7fa613666c234e1090e91427314ee18ec6420e9426cf4e7f314056813", digest)
}

func TestImageRegistryRewrite(t *testing.T) {
	cfg := &config.Config{
		RegistryRewrites: map[string]string{
			"docker.io": "registry-1.docker.io",
		},
		EnableRewrite: true,
	}
	di := NewImage(cfg, "docker.io/library", "alpine", "3.12")
	assert.Equal(t, "registry-1.docker.io/library/alpine", di.RegistryRepository())
}
