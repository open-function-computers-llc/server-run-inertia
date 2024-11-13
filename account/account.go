package account

type Account struct {
	IsLocked                  bool     `json:"isLocked"`
	Name                      string   `json:"name"`
	UptimeURI                 string   `json:"uptimeURI"`
	Username                  string   `json:"username"`
	PrimaryDomain             string   `json:"domain"`
	AlternateDomains          []string `json:"alternateDomains"`
	AlwaysUnlockedDirectories []string `json:"alwaysUnlockedDirectories"`
	PubKey                    string   `json:"sshPubKey"`
	CreatedAt                 string   `json:"createdAt"`
}

func New(a string) (Account, error) {
	return Account{
		Name: a,
	}, nil
}

func (a *Account) LoadStatus() error {
	return nil
}
