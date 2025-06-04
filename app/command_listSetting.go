package app

import (
	"errors"
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/urfave/cli/v2"
)

func listSetting() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Display some information from the account settings file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "account",
				Usage:    "Account you're retrieving data from",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "setting",
				Usage:    "The setting you are tring to see (databases|alternate-domains|is-locked)",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			a := cCtx.String("account")
			account, err := account.Find(a)
			if err != nil {
				return err
			}

			s := cCtx.String("setting")

			if s == "databases" {
				for _, db := range account.Databases {
					fmt.Println(db.Name)
				}
				return nil
			}
			if s == "alternate-domains" {
				for _, d := range account.AlternateDomains {
					fmt.Println(d)
				}
				return nil
			}
			if s == "is-locked" {
				if account.IsLocked {
					fmt.Println("true")
					return nil
				}
				fmt.Println("false")
				return nil
			}
			if s == "domain" || s == "primary-domain" {
				fmt.Println(account.PrimaryDomain)
				return nil
			}

			return errors.New("WARNING: invalid setting requested")
		},
	}
}
