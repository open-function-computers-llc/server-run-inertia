package client

import (
	"os"
	"sync"

	gt "code.gitea.io/sdk/gitea"
)

var client *gt.Client
var once sync.Once
var err error

func Client() (*gt.Client, error) {
	once.Do(func() {
		client, err = gt.NewClient(os.Getenv("GITEA_URL"), gt.SetToken(os.Getenv("GITEA_TOKEN")))

		if err != nil {
			return
		}
	})

	return client, err
}
