package main

import (
	"embed"
	"errors"
	"os"
	"strconv"

	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
	"github.com/open-function-computers-llc/server-run-inertia/server"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed dist/*
var dist embed.FS

func main() {
	port, url, err := verifyValidENV()
	if err != nil {
		panic("Trouble setting up server: " + err.Error())
	}

	err = scriptrunner.VerifyAllScriptsExist()
	if err != nil {
		panic("Missing script: " + err.Error())
	}

	s, err := server.New(port, url, dist)
	if err != nil {
		panic("Trouble setting up server: " + err.Error())
	}

	err = s.Serve()
	if err != nil {
		panic("Trouble serving from server: " + err.Error())
	}
}

func verifyValidENV() (int, string, error) {
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
