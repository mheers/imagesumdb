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
	// registry of the image
	registry string
	// repository of the image (=name)
	repository string
	// tag of the image
	tag string
	// wheather the image is a usual image or a general oci artifact
	oci bool
	// cfg is the config
	cfg       *config.Config
	vulncheck *Vulncheck
}

func NewImage(cfg *config.Config, registry, repository, tag string) *Image {
	return &Image{
		registry:   registry,
		repository: repository,
		tag:        tag,
		cfg:        cfg,
		oci:        false,
	}
}

func NewOCIImage(cfg *config.Config, registry, repository, tag string) *Image {
	return &Image{
		registry:   registry,
		repository: repository,
		tag:        tag,
		cfg:        cfg,
		oci:        true,
	}
}

func (i *Image) RegistryRepositoryPlain() string {
	r := i.registry
	if r != "" {
		r += "/"
	}
	return fmt.Sprintf("%s%s", r, i.repository)
}

func (i *Image) RegistryRepository() string {
	r := i.Registry()
	if r != "" {
		r += "/"
	}
	return fmt.Sprintf("%s%s", r, i.Repository())
}

func (i *Image) RegistryRepositoryTagPlain() string {
	return fmt.Sprintf("%s:%s", i.RegistryRepositoryPlain(), i.tag)
}

func (i *Image) RegistryRepositoryTag() string {
	return fmt.Sprintf("%s:%s", i.RegistryRepository(), i.Tag())
}

func (i *Image) RegistryRepositoryTagDigestPlain() string {
	return fmt.Sprintf("%s:%s@%s", i.RegistryRepositoryPlain(), i.TagPlain(), i.MustDigest())
}

func (i *Image) RegistryRepositoryTagDigest() string {
	return fmt.Sprintf("%s:%s@%s", i.RegistryRepository(), i.TagPlain(), i.MustDigest())
}

func (i *Image) StringPlain() string {
	return i.RegistryRepositoryTagPlain()
}

func (i *Image) String() string {
	return i.RegistryRepositoryTag()
}

func (i *Image) RegistryPlain() string {
	return i.registry
}

func (i *Image) Registry() string {
	return imageRegistryRewrite(i.cfg, i.RegistryPlain())
}

func (i *Image) Repository() string {
	return i.repository
}

func (i *Image) TagPlain() string {
	return i.tag
}

func (i *Image) Tag() string {
	if i.cfg == nil {
		return i.tag
	}
	if i.cfg.ForceDigest {
		digest := i.MustDigest()
		return fmt.Sprintf("%s@%s", i.tag, digest)
	}
	return i.tag
}

func (i *Image) imageCloser() (types.ImageCloser, error) {
	ref, err := docker.ParseReference(fmt.Sprintf("//%s", i.RegistryRepositoryTagPlain()))
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

func (i *Image) MustDigest() string {
	digest, err := i.Digest()
	if err != nil {
		panic(err)
	}
	return digest
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
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
		Config        struct {
			MediaType string `json:"mediaType"`
			Size      int    `json:"size"`
			Digest    string `json:"digest"`
		} `json:"config"`
		Layers []struct {
			MediaType string `json:"mediaType"`
			Size      int    `json:"size"`
			Digest    string `json:"digest"`
		} `json:"layers"`
		Manifests []struct {
			Digest string `json:"digest"`
		} `json:"manifests"`
	}

	var m manifest
	err = json.Unmarshal([]byte(manifestS), &m)
	if err != nil {
		return "", err
	}

	var digest string
	if len(m.Manifests) == 0 {
		digest = m.Config.Digest
	} else {
		digest = m.Manifests[0].Digest
	}

	return digest, nil
}

func (i *Image) IsOCI() bool {
	return i.oci
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
