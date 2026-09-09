package app

import (
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/open-function-computers-llc/server-run-inertia/errx"
	"github.com/urfave/cli/v2"
)

func envExport() *cli.Command {
	return &cli.Command{
		Name:    "set-env",
		Aliases: []string{"env"},
		Usage:   "Populate the system ENV for use in scripts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "account",
				Usage:    "Account you're retrieving data from",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "database",
				Usage: "Populate DBUSER DBHOST DBPASSWORD DBNAME for a given database",
			},
		},
		Action: func(cCtx *cli.Context) error {
			a := cCtx.String("account")
			account, err := account.Find(a)
			if err != nil {
				return err
			}

			if cCtx.String("database") != "" {
				for _, db := range account.Databases {
					if db.Name == cCtx.String("database") {
						fmt.Println("export DB_ACCOUNT='" + account.Name + "'")
						fmt.Println("export DB_USER='" + db.Username + "'")
						fmt.Println("export DB_NAME='" + db.Name + "'")
						fmt.Println("export DB_PASSWORD='" + db.Password + "'")
						fmt.Println("export DB_HOST='" + db.Host + "'")
						return nil
					}
				}
				return errx.Wrapf("database not found in account", fmt.Errorf("database %q not found in account %q", cCtx.String("database"), account.Name), ErrorDatabaseNotFoundInAccount)
			}

			return errx.Wrapf("invalid param provided", fmt.Errorf("database is missing"), ErrorInvalidParam)
		},
	}
}
