package vcs

import "time"

// (not in use yet) - should implement
type Repository struct {
	Name      string
	Url       string
	HTMLUrl   string
	Private   bool
	CreatedAt time.Time

	Id string
}

type DeployKey struct {
	Name      string
	Key       string
	ReadOnly  bool
	CreatedAt time.Time
	CreatedBy string

	Id string
}

type Provider interface {
	CreateRepository(name string) error
	ListRepositories() ([]string, error)
	ListDeployKeys(repository string) ([]DeployKey, error)
	AddDeployKey(repository string, key string, name string, readOnly bool) (DeployKey, error)
	RemoveDeployKey(repository string, id string) error
}
