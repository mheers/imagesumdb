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
		Registry:   i.registry,
		Repository: i.repository,
		Tag:        i.tag,
		OCI:        i.oci,
		Digest:     digest,
		Vulncheck:  i.vulncheck,
	}, nil
}

func (i *ImagePersistance) ToImage() *Image {
	return &Image{
		registry:   i.Registry,
		repository: i.Repository,
		tag:        i.Tag,
		oci:        i.OCI,
		vulncheck:  i.Vulncheck,
	}
}
