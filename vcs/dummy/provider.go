package dummy

type dummyProvider struct {
}

func NewProvider() (dummyProvider, error) {
	p := dummyProvider{}
	return p, nil
}

func (p dummyProvider) ListRepositories() ([]string, error) {
	return []string{""}, nil
}
func (p dummyProvider) CreateRepository(name string) error {
	return nil
}
