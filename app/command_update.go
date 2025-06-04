package app

import (
	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/urfave/cli/v2"
)

type updateableProperty struct {
	name     string
	usage    string
	callback func(cCtx *cli.Context) error
}

var updateableProperties = []updateableProperty{
	{
		name:     "locked",
		usage:    "Set locked status for this account (true|false)",
		callback: setLock,
	},
	{
		name:     "set-domain",
		usage:    "Update the main domain for the selected account",
		callback: setPrimaryDomain,
	},
	{
		name:     "add-domain",
		usage:    "Add a domain to the selected account",
		callback: addAlternateDomain,
	},
	{
		name:     "remove-domain",
		usage:    "Remove a domain from the selected account",
		callback: removeAlternateDomain,
	},
	{
		name:     "add-database",
		usage:    "Add a database to the selected account",
		callback: addDatabase,
	},
	{
		name:  "db-user",
		usage: "Username when calling \"add-database\"",
	},
	{
		name:  "db-name",
		usage: "Database name when calling \"add-database\"",
	},
	{
		name:  "db-host",
		usage: "Database host when calling \"add-database\"",
	},
	{
		name:  "db-password",
		usage: "Password when calling \"add-database\"",
	},
}

func update() *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "account",
			Usage:    "Account you're updating manually (required)",
			Required: true,
		},
	}

	for _, f := range updateableProperties {
		flags = append(flags, &cli.StringFlag{
			Name:  f.name,
			Usage: f.usage,
		})
	}

	return &cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update a specific account's settings file",
		Flags:   flags,
		Action:  updateAction,
	}
}

func updateAction(cCtx *cli.Context) error {
	accountName := cCtx.String("account")

	_, err := account.Find(accountName)
	if err != nil {
		return err
	}

	for _, property := range updateableProperties {
		val := cCtx.String(property.name)
		if val != "" && property.callback != nil {
			return property.callback(cCtx)
		}
	}

	return &ErrorInvalidArguments{}
}

func setLock(cCtx *cli.Context) error {
	a, _ := account.Find(cCtx.String("account"))
	if cCtx.String("locked") == "true" {
		a.Lock()
		return nil
	}
	if cCtx.String("locked") == "false" {
		a.Unlock()
		return nil
	}
	return &ErrorInvalidArguments{}
}

func setPrimaryDomain(cCtx *cli.Context) error {
	a, _ := account.Find(cCtx.String("account"))
	return a.SetPrimaryDomain(cCtx.String("set-domain"))
}

func addAlternateDomain(cCtx *cli.Context) error {
	a, _ := account.Find(cCtx.String("account"))

	return a.AddAlternateDomain(cCtx.String("add-domain"))
}

func removeAlternateDomain(cCtx *cli.Context) error {
	a, _ := account.Find(cCtx.String("account"))

	return a.RemoveAlternateDomain(cCtx.String("remove-domain"))
}

func addDatabase(cCtx *cli.Context) error {
	a, _ := account.Find(cCtx.String("account"))

	return a.AddDatabase(cCtx.String("db-host"), cCtx.String("db-user"), cCtx.String("db-password"), cCtx.String("db-name"))
}
