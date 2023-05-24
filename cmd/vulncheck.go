package main

import (
	"github.com/mheers/imagesumdb/config"
	"github.com/mheers/imagesumdb/db"
	"github.com/mheers/imagesumdb/image"
	"github.com/mheers/imagesumdb/vulncheck"
	"github.com/spf13/cobra"
)

var (
	dbFile       string // dbFile is the path to the db file
	vulncheckCmd = &cobra.Command{
		Use:   "check",
		Short: "check for vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.Config{
				DBFile: dbFile,
			}
			d := db.NewDB(cfg)
			if err := d.Read(); err != nil {
				return err
			}
			if err := checkVulnerabilities(d); err != nil {
				return err
			}
			if err := d.Write(); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	vulncheckCmd.Flags().StringVarP(&dbFile, "db", "d", "db.yaml", "path to db file")
	rootCmd.AddCommand(vulncheckCmd)
}

func checkVulnerabilities(db *db.DB) error {
	// run vulncheck on all images in db.Images
	for _, img := range db.GetAll() {
		if img.IsOCI() {
			continue
		}
		report, err := vulncheck.Scan(img)
		if err != nil {
			return err
		}

		img.SetVulncheck(image.NewVulncheck(report))
	}
	return nil
}
