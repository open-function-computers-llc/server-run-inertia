package github

type gitHubProvider struct {
	API_TOKEN string
}

func NewProvider() (gitHubProvider, error) {
	p := gitHubProvider{}
	return p, nil
}
