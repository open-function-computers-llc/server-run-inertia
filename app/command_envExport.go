package app

import (
	"errors"
	"fmt"

	"github.com/open-function-computers-llc/server-run-inertia/account"
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
				return errors.New("Couldn't find database " + cCtx.String("database") + " in account " + account.Name)
			}

			return errors.New("WARNING: no valid flags set!")
		},
	}
}
