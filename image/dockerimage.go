package image

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
	"github.com/mheers/imagesumdb/config"
)

type Image struct {
	// registry with repository of the image
	repository string
	// tag of the image
	tag string
	// cfg is the config
	cfg       *config.Config
	vulncheck *Vulncheck
}

func NewImage(cfg *config.Config, repository, tag string) *Image {
	return &Image{
		repository: repository,
		tag:        tag,
		cfg:        cfg,
	}
}

func (i *Image) RepoTag() string {
	return i.repository + ":" + i.tag
}

func (i *Image) String() string {
	return i.Repository() + ":" + i.Tag()
}

func (i *Image) Repository() string {
	if i.cfg == nil {
		return i.repository
	}
	return imageRegistryRewrite(i.cfg, i.repository)
}

func (i *Image) Tag() string {
	if i.cfg == nil {
		return i.tag
	}
	if i.cfg.ForceDigest {
		digest, err := i.Digest()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s@%s", i.tag, digest)
	}
	return i.tag
}

func (i *Image) imageCloser() (types.ImageCloser, error) {
	ref, err := docker.ParseReference(fmt.Sprintf("//%s", i.RepoTag()))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	img, err := ref.NewImage(ctx, nil)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// func (i *Image) pull()        {}
// func (i *Image) writeDigest() {}

func (i *Image) Manifest() (string, error) {
	img, err := i.imageCloser()
	if err != nil {
		return "", err
	}
	defer img.Close()

	ctx := context.Background()
	b, _, err := img.Manifest(ctx)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (i *Image) Digest() (string, error) {
	img, err := i.imageCloser()
	if err != nil {
		return "", err
	}
	defer img.Close()

	manifestS, err := i.Manifest()
	if err != nil {
		return "", err
	}

	type manifest struct {
		Manifests []struct {
			Digest string `json:"digest"`
		} `json:"manifests"`
	}

	var m manifest
	err = json.Unmarshal([]byte(manifestS), &m)
	if err != nil {
		return "", err
	}

	digest := m.Manifests[0].Digest

	return digest, nil
}

func imageRegistryRewrite(cfg *config.Config, src string) string {
	if !cfg.EnableRewrite || len(cfg.RegistryRewrites) == 0 {
		return src
	}

	prefix := strings.Split(src, "/")[0]

	if rewrite, ok := cfg.RegistryRewrites[prefix]; ok {
		return strings.Replace(src, prefix, rewrite, 1)
	}
	return src
}
