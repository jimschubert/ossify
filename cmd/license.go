package cmd

import (
	"fmt"
	"github.com/jimschubert/ossify/config"
	"github.com/jimschubert/ossify/licenses"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var licenseFlags *LicenseFlags

type LicenseFlags struct {
	licenseId       string
	licenseTemplate string
	keyword         []string
	search          string
	details         bool
}

func init() {
	licenseFlags = &LicenseFlags{}
	rootCmd.AddCommand(licenseCmd)

	// license
	licenseCmd.AddCommand(listLicenseCmd)
	licenseCmd.AddCommand(addLicenseCmd)

	licenseCmd.Flags().StringVarP(&licenseFlags.licenseId, "id", "i", "",
		"Get details about a single license by ID.")
	licenseCmd.Flags().StringSliceVar(&licenseFlags.keyword, "keyword", []string{},
		"Keywords to filter remote licenses by\n\t(copyleft,discouraged,international,miscellaneous,\n\t non-reusable,obsolete,osi-approved,permissive,\n\t popular,redundant,retired,special-purpose)")
	licenseCmd.Flags().StringVar(&licenseFlags.search, "search", "",
		"Search term to query across all known license metadata.")
	licenseCmd.Flags().BoolVar(&licenseFlags.details, "details", false,
		"When included with the id option, prints only the details of the requested license rather than the license text.")

	// license add
	addLicenseCmd.Flags().StringVarP(&licenseFlags.licenseId, "id", "i", "",
		"The identifier to be associated with your customized license. This will take "+
			"precedence over 'public' ids.")

	addLicenseCmd.Flags().StringVarP(&licenseFlags.licenseTemplate, "template", "t", "",
		"The template to add for the given identifier.")

	// license list
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Manage open-source licenses",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ConfigManager.Load()
		failOnError(err)

		// TODO: define _where_ local configuration will be held (e.g. ~/.config/ossify), used by this and Add
		//		 then, pull from OSI list, and merge our local licenses on top of that.

		allLicenses, err := licenses.LoadLicenses()
		failOnError(err)

		id := licenseFlags.licenseId
		keywords := licenseFlags.keyword
		search := licenseFlags.search

		if len(id) == 0 && len(keywords) == 0 && len(search) == 0 {
			keywords = append(keywords, "popular")
		}

		if len(id) > 0 {
			license := allLicenses.FindById(id)
			details := licenseFlags.details
			if license != nil {
				if details {
					_ = license.Print()
				} else {
					err := licenses.PrintLicenseText(license.Id, conf.LicensePath)
					failOnError(err)
				}
			}
		} else if len(keywords) > 0 {
			for _, keyword := range keywords {
				keywordLicenses := allLicenses.FindByKeyword(keyword)
				for _, byKeyword := range *keywordLicenses {
					_ = byKeyword.Print()
				}
			}
		} else if len(search) > 0 {
			searchResults := allLicenses.Search(search)
			for _, result := range *searchResults {
				_ = result.Print()
			}
		}
	},
}

var addLicenseCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new custom license (local-only) to the list of known licenses.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		config, err := config.ConfigManager.Load()
		failOnError(err)

		licensePath := config.LicensePath
		if licensePath == "" {
			err = errors.New("invalid license path: please update your configuration and try again")
			failOnError(err)
		}

		err = os.MkdirAll(licensePath, 0700)
		failOnError(err)

		if licenseFlags.licenseTemplate == "" {
			licenseFlags.licenseTemplate = args[0]

			if licenseFlags.licenseTemplate == "" {
				err = errors.New("invalid template: you must provide a template value")
				failOnError(err)
			}
		}

		if licenseFlags.licenseId == "" {
			err = errors.New("invalid id: you must provide a id value")
			failOnError(err)
		}

		data, _ := ioutil.ReadFile(licenseFlags.licenseTemplate)

		err = os.MkdirAll(licensePath, 0700)
		failOnError(err)

		targetFile := path.Join(licensePath, licenseFlags.licenseId)

		// TODO: Document how this allows users to specify default text for a license
		err = ioutil.WriteFile(targetFile, data, 0644)
		failOnError(err)

		log.Println(fmt.Sprintf("Saved license with id %s to %s", licenseFlags.licenseId, targetFile))
	},
}

var listLicenseCmd = &cobra.Command{
	Use:   "list",
	Short: "Presents a list of known licenses.",
	Run: func(cmd *cobra.Command, args []string) {
		licenses, err := licenses.LoadLicenses()
		failOnError(err)

		for _, license := range *licenses {
			err = license.Print()
			failOnError(err)
		}
	},
}
