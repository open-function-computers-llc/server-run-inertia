package gitea

type gitTeaProvider struct {
	AUTH_KEY string
	URL      string
}

func NewProvider() (gitTeaProvider, error) {
	p := gitTeaProvider{}
	return p, nil
}
