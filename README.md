# ImageSumDB

> A simple container image database for vuln scanning and version pinning

ImageSumDB provides a flexible database as a dependency for projects that deploy dockerimages. It is a simple database that stores the name of the image and the sum of the image. The sum is calculated by the sha256sum of the image. The database is stored in a file called `db.yaml` in the current working directory.

## Features
- Add images to the database
- Calculate the sum of the image
- Check if the image is in the database
- Vulnerability check of the image

## Usage
See godoc for more information.

## Example

```go
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

err = dbInstance.Set("baseimage", image.NewImage(&cfg, "alpine", "3.16"))
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
```

# Alternatives

- https://github.com/safe-waters/docker-lock
