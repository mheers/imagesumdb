package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/mheers/imagesumdb/config"
	"github.com/mheers/imagesumdb/image"
	"github.com/mheers/imagesumdb/pluginregistry"
	"gopkg.in/yaml.v3"
)

type DB struct {
	// Cfg is the config
	Cfg              *config.Config
	Images           map[string]*image.Image
	imagePersistance map[string]*image.ImagePersistance
}

func NewDB(cfg *config.Config) *DB {
	return &DB{
		Cfg:    cfg,
		Images: make(map[string]*image.Image),
	}
}

func (db *DB) Get(name string) (*image.Image, error) {
	// get image from db.Images
	return db.Images[name], nil
}

func (db *DB) GetAll() map[string]*image.Image {
	return db.Images
}

func (db *DB) AddOCI(name, registry, repository, tag string) error {
	return db.add(name, image.NewOCIImage(db.Cfg, registry, repository, tag))
}

func (db *DB) Add(name, registry, repository, tag string) error {
	return db.add(name, image.NewImage(db.Cfg, registry, repository, tag))
}

func (db *DB) add(name string, img *image.Image) error {
	// set image in db.Images
	if db.Images[name] != nil {
		return fmt.Errorf("image for %s already exists", name)
	}
	db.Images[name] = img
	return nil
}

func (db *DB) Import(db2 *DB) error {
	// import db2.Images to db.Images
	for name, img := range db2.Images {
		db.Images[name] = img
	}
	return nil
}

func (db *DB) Set(name string, img *image.Image) error {
	// set image in db.Images
	db.Images[name] = img
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
	imagePersistance, _, err := GetImagesFromFile(db.Cfg.DBFile)
	if err != nil {
		return err
	}
	db.imagePersistance = imagePersistance

	return nil
}

func (db *DB) ToImagePersistance() (map[string]*image.ImagePersistance, error) {
	if db.imagePersistance == nil {
		err := db.ReadImagePersistance()
		if err != nil {
			return nil, err
		}
	}
	return db.imagePersistance, nil
}

func (db *DB) Read() error {
	_, images, err := GetImagesFromFile(db.Cfg.DBFile)
	if err != nil {
		return err
	}
	db.Images = images

	return nil
}

func getDir(path string) string {
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return "./"
	}
	return path[:index]
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
	return createFileWithDirs(db.Cfg.DBFile)
}

func (db *DB) Write() error {
	// create db.cfg.DBFile if it does not exist
	// write db.Images to db.cfg.DBFile

	result, err := db.getPersistanceMap()
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return nil
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
	vchecker := pluginregistry.GetVulncheck()
	for _, img := range db.GetAll() {
		if img.IsOCI() {
			continue
		}
		report, err := vchecker.Scan(img)
		if err != nil {
			return err
		}

		img.SetVulncheck(image.NewVulncheck(report))
	}
	return nil
}

func (db *DB) CompareSetImagesWithPersistance() error {
	return db.CompareImages(db.Images)
}

func (db *DB) CompareImages(new map[string]*image.Image) error {
	old := db.imagePersistance
	// check if hashes have changed
	for name, img := range new {
		if old[name] == nil {
			fmt.Printf("new image found for %s: %s\n", name, img.RegistryRepositoryTag())
		} else {
			newDigest, err := img.Digest()
			if err != nil {
				return err
			}
			if old[name].Digest != newDigest {
				fmt.Printf("image changed for %s: %s\n", name, img.RegistryRepositoryTag())
			}
		}
	}
	return nil
}

func (db *DB) getPersistanceMap() (map[string]*image.ImagePersistance, error) {
	// convert db.Images to map[string]*image.ImagePersistance
	result := make(map[string]*image.ImagePersistance)
	for name, img := range db.Images {
		p, err := img.ToPersistance(name)
		if err != nil {
			return nil, err
		}
		result[name] = p
	}
	return result, nil
}
