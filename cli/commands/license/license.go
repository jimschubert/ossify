package license

import (
	"github.com/jimschubert/ossify/config"
	licenseUtil "github.com/jimschubert/ossify/licenses"
	"gopkg.in/urfave/cli.v1"
)

func Command() cli.Command {
	return cli.Command{
		Name:     "license",
		Usage:    "Manage open-source licenses.",
		HelpName: "license",
		Category: "Manage",
		Subcommands: []cli.Command{
			Add,
			List,
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "id",
				Usage: "Get details about a single license by ID.",
			},
			// list of options can be found with:
			// jq -r '.[].keywords | @tsv ' licenses.json  | tr '\t' '\n' | sort | uniq
			cli.StringSliceFlag{
				Name:  "keyword",
				Usage: "Keywords to filter remote licenses by\n\t(copyleft,discouraged,international,miscellaneous,\n\t non-reusable,obsolete,osi-approved,permissive,\n\t popular,redundant,retired,special-purpose)",
				Value: &cli.StringSlice{},
			},
			cli.StringFlag{
				Name: "search",
				Usage: "Search term to query across all known license metadata.",
			},
			cli.BoolFlag{
				Name: "details",
				Usage: "When included with `id`, prints only the details of the requested license rather than the license text.",
			},
		},
		Action: func(c *cli.Context) error {
			conf, err := config.ConfigManager.Load()
			if err != nil {
				return err
			}
			// TODO: define _where_ local configuration will be held (e.g. ~/.config/ossify), used by this and Add
			//		 then, pull from OSI list, and merge our local licenses on top of that.

			licenses, err := licenseUtil.LoadLicenses()
			if err != nil {
				return err
			}

			id := c.String("id")
			keywords := c.StringSlice("keyword")
			search := c.String("search")

			if len(id) == 0 && len(keywords) == 0 && len(search) == 0 {
				keywords = append(keywords, "popular")
			}

			if len(id) > 0 {
				license := licenses.FindById(id)
				details := c.Bool("details")
				if license != nil {
					if details {
						_ = license.Print()
					} else {
						err := licenseUtil.PrintLicenseText(license.Id, conf.LicensePath)
						if err != nil {
							return err
						}
					}
				}
			} else if len(keywords) > 0 {
				for _, keyword := range keywords {
					keywordLicenses := licenses.FindByKeyword(keyword)
					for _, byKeyword := range *keywordLicenses {
						_ = byKeyword.Print()
					}
				}
			} else if len(search) > 0 {
				searchResults := licenses.Search(search)
				for _, result := range *searchResults {
					_ = result.Print()
				}
			}

			return nil
		},
	}
}
