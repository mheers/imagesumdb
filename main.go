package main

import (
	"github.com/mheers/imagesumdb/config"
	"github.com/mheers/imagesumdb/db"
	"github.com/mheers/imagesumdb/image"
)

func main() {
	cfg := config.Config{
		RegistryRewrites: map[string]string{
			"docker.io": "registry-1.docker.io",
		},
		EnableRewrite: true,
		DBFile:        "db.yaml",
		ForceDigest:   true,
	}

	dbInstance := db.NewDB(&cfg)
	err := dbInstance.ReadImagePersistance()
	if err != nil {
		panic(err)
	}

	err = dbInstance.Set("test", image.NewImage(&cfg, "alpine", "3.11"))
	if err != nil {
		panic(err)
	}

	err = dbInstance.Set("newer", image.NewImage(&cfg, "alpine", "3.16"))
	if err != nil {
		panic(err)
	}

	err = dbInstance.CompareSetImagesWithPersistance()
	if err != nil {
		panic(err)
	}

	err = dbInstance.Vulncheck()
	if err != nil {
		panic(err)
	}

	err = dbInstance.Write()
	if err != nil {
		panic(err)
	}
}
