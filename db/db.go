package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/mheers/imagesumdb/config"
	"github.com/mheers/imagesumdb/image"
	"gopkg.in/yaml.v3"
)

type DB struct {
	// cfg is the config
	cfg              *config.Config
	images           map[string]*image.Image
	imagePersistance map[string]*image.ImagePersistance
}

func NewDB(cfg *config.Config) *DB {
	return &DB{
		cfg:    cfg,
		images: make(map[string]*image.Image),
	}
}

func (db *DB) Get(name string) (*image.Image, error) {
	// get image from db.Images
	return db.images[name], nil
}

func (db *DB) Add(name string, repository, tag string) error {
	// set image in db.Images
	if db.images[name] != nil {
		return fmt.Errorf("image for %s already exists", name)
	}
	db.images[name] = image.NewImage(db.cfg, repository, tag)
	return nil
}

func (db *DB) Set(name string, img *image.Image) error {
	// set image in db.Images
	db.images[name] = img
	return nil
}

func GetImagesFromFile(path string) (map[string]*image.ImagePersistance, map[string]*image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = createFileWithDirs(path)
			if err != nil {
				return nil, nil, err
			}
			f, err = os.Open(path)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	result := make(map[string]*image.ImagePersistance)
	err = yaml.NewDecoder(f).Decode(&result)
	if err != nil {
		// checks if EOF
		if err.Error() == "EOF" {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// convert map[string]*image.ImagePersistance to db.Images
	images := make(map[string]*image.Image)
	for name, img := range result {
		images[name] = img.ToImage()
	}
	return result, images, nil
}

func (db *DB) ReadImagePersistance() error {
	imagePersistance, _, err := GetImagesFromFile(db.cfg.DBFile)
	if err != nil {
		return err
	}
	db.imagePersistance = imagePersistance

	return nil
}

func (db *DB) Read() error {
	_, images, err := GetImagesFromFile(db.cfg.DBFile)
	if err != nil {
		return err
	}
	db.images = images

	return nil
}

func getDir(path string) string {
	return path[:strings.LastIndex(path, "/")]
}

func createFileWithDirs(path string) (*os.File, error) {
	// create file with dirs
	err := os.MkdirAll(getDir(path), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (db *DB) create() (*os.File, error) {
	return createFileWithDirs(db.cfg.DBFile)
}

func (db *DB) Write() error {
	// create db.cfg.DBFile if it does not exist
	// write db.Images to db.cfg.DBFile

	result, err := db.getPersistanceMap()
	if err != nil {
		return err
	}

	f, err := db.create()
	if err != nil {
		return err
	}

	resultString, err := yaml.Marshal(result)
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(resultString))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Vulncheck() error {
	// run vulncheck on all images in db.Images
	for _, img := range db.images {
		_, err := img.Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) CompareSetImagesWithPersistance() error {
	return db.CompareImages(db.images)
}

func (db *DB) CompareImages(new map[string]*image.Image) error {
	old := db.imagePersistance
	// check if hashes have changed
	for name, img := range new {
		if old[name] == nil {
			fmt.Printf("new image found for %s: %s\n", name, img.RepoTag())
		} else {
			newDigest, err := img.Digest()
			if err != nil {
				return err
			}
			if old[name].Digest != newDigest {
				fmt.Printf("image changed for %s: %s\n", name, img.RepoTag())
			}
		}
	}
	return nil
}

func (db *DB) getPersistanceMap() (map[string]*image.ImagePersistance, error) {
	// convert db.Images to map[string]*image.ImagePersistance
	result := make(map[string]*image.ImagePersistance)
	for name, img := range db.images {
		p, err := img.ToPersistance(name)
		if err != nil {
			return nil, err
		}
		result[name] = p
	}
	return result, nil
}
