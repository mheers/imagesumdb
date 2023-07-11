package image

type ImagePersistance struct {
	Name       string     `yaml:"name"`
	Registry   string     `yaml:"registry"`
	Repository string     `yaml:"repository"`
	Tag        string     `yaml:"tag"`
	OCI        bool       `yaml:"oci"`
	Digest     string     `yaml:"digest"`
	Vulncheck  *Vulncheck `yaml:"vulncheck"`
}

func (i *Image) ToPersistance(name string) (*ImagePersistance, error) {
	digest, err := i.Digest()
	if err != nil {
		return nil, err
	}
	return &ImagePersistance{
		Name:       name,
		Registry:   i.I_registry,
		Repository: i.I_repository,
		Tag:        i.I_tag,
		OCI:        i.I_oci,
		Digest:     digest,
		Vulncheck:  i.I_vulncheck,
	}, nil
}

func (i *ImagePersistance) ToImage() *Image {
	return &Image{
		I_registry:   i.Registry,
		I_repository: i.Repository,
		I_tag:        i.Tag,
		I_oci:        i.OCI,
		I_vulncheck:  i.Vulncheck,
	}
}
