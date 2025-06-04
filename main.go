package main

import (
	"embed"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/open-function-computers-llc/server-run-inertia/app"
	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
	"github.com/urfave/cli/v2"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed dist/*
var dist embed.FS

//go:embed assets/logo-default.png
var defaultLogo []byte

// version can be changed by build flags
var Version = "latest-dev"

func main() {
	port, url, err := verifyValidENV()
	if err != nil {
		panic("Trouble setting up server: " + err.Error())
	}

	err = scriptrunner.VerifyAllScriptsExist()
	if err != nil {
		panic(err.Error())
	}

	app := &cli.App{
		Name:     "Server Run",
		Usage:    "Manage this web server",
		Version:  Version,
		Commands: app.AvailableCommands(port, url, dist, defaultLogo),
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func verifyValidENV() (int, string, error) {
	// check global env file
	_, err := os.Stat("/etc/server-run.env")
	if err == nil {
		// global env file exists! let's load it
		godotenv.Load("/etc/server-run.env")
	}

	requiredENV := []string{
		"APP_PORT",
		"APP_URL",
		"AUTH_USER",
		"AUTH_PASSWORD",
		"UPTIME_USER",
		"UPTIME_PASS",
		"UPTIME_DOMAIN",
		"THIRD_PARTY_ACCESS_TOKEN",
		"WEBSITES_ROOT",
		"ACCOUNTS_ROOT",
	}

	for _, env := range requiredENV {
		check := os.Getenv(env)
		if check == "" {
			return 0, "", errors.New("missing env: " + env)
		}
	}

	webPort := os.Getenv("APP_PORT")
	portInt, err := strconv.Atoi(webPort)
	if err != nil {
		panic("Invalid APP_PORT: " + err.Error())
	}

	return portInt, os.Getenv("APP_URL"), nil
}
