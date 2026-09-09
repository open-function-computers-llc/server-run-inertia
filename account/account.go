package account

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/open-function-computers-llc/server-run-inertia/errx"
)

type Account struct {
	IsLocked                  bool       `json:"isLocked"`
	IsLockedDownDisabled      bool       `json:"isLockedDownDisabled"`
	Name                      string     `json:"name"`
	UptimeURL                 string     `json:"uptimeURL"`
	Username                  string     `json:"username"`
	PrimaryDomain             string     `json:"domain"`
	AlternateDomains          []string   `json:"alternateDomains"`
	AlwaysUnlockedDirectories []string   `json:"alwaysUnlockedDirectories"`
	Databases                 []Database `json:"databases"`
	PubKey                    string     `json:"sshPubKey"`
	CreatedAt                 time.Time  `json:"createdAt"`
	UpdatedAt                 time.Time  `json:"updatedAt"`
}

func Find(name string) (*Account, error) {
	a := Account{
		Name: name,
	}

	err := a.refreshStatus()
	if err != nil {
		return nil, err
	}

	if a.CreatedAt.IsZero() {
		folderPath, _ := a.stateDirectory()
		fmt.Println("folderpath: " + folderPath)
		if err != nil {
			return nil, err
		}

		folder, err := os.Stat(folderPath)
		if err != nil {
			return nil, err
		}

		a.CreatedAt = folder.ModTime()
	}

	return &a, err
}

func (a *Account) stateDirectory() (string, error) {
	if strings.TrimSpace(a.Name) == "" {
		return "", ErrorInvalidAccountName
	}

	accountsRoot := os.Getenv("ACCOUNTS_ROOT")
	if accountsRoot == "" {
		accountsRoot = "/tmp/"
	}

	accountHome := accountsRoot + a.Name
	_, err := os.Stat(accountHome)
	if err == nil {
		return accountHome, nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(accountHome, 0755)
		if err != nil {
			return "", err
		}
	}
	fmt.Println("account home: " + accountHome)
	return accountHome, nil
}

func (a *Account) refreshStatus() error {
	// account "home" must exist
	folderPath, err := a.stateDirectory()
	if err != nil {
		return err
	}

	settingsFile := folderPath + "/settings.json"
	_, err = os.Stat(settingsFile)
	if err != nil {
		folder, err := os.Stat(folderPath)
		if err != nil {
			return err
		}
		a.CreatedAt = folder.ModTime()
		a.writeSettingsFile()
	}

	status, err := os.ReadFile(settingsFile)
	if err != nil {
		return errx.Wrapf("reading settings file", err, ErrorReadingStateFile)
	}

	err = json.Unmarshal(status, a)
	if err != nil {
		return errx.Wrapf("parsing settings file", err, ErrorReadingStateFile)
	}

	return nil
}

func (a *Account) writeSettingsFile() error {
	if a.AlternateDomains == nil {
		a.AlternateDomains = []string{}
	}
	if a.AlwaysUnlockedDirectories == nil {
		a.AlwaysUnlockedDirectories = []string{}
	}
	if a.Databases == nil {
		a.Databases = []Database{}
	}

	a.UpdatedAt = time.Now()

	json, _ := json.Marshal(a)
	folderPath, _ := a.stateDirectory()
	return os.WriteFile(folderPath+"/settings.json", json, 0644)
}

func (a *Account) Lock() error {
	a.IsLocked = true
	return a.writeSettingsFile()
}

func (a *Account) Unlock() error {
	a.IsLocked = false
	return a.writeSettingsFile()
}

func (a *Account) SetLockedDownDisabled(enabled bool) error {
	a.IsLockedDownDisabled = enabled
	return a.writeSettingsFile()
}

func (a *Account) SetPrimaryDomain(domain string) error {
	if strings.TrimSpace(domain) == "" {
		return nil
	}
	// if strings.TrimSpace(a.PrimaryDomain) != "" {
	// 	a.AlternateDomains = append(a.AlternateDomains, a.PrimaryDomain)
	// }

	if slices.Contains(a.AlternateDomains, domain) {
		newList := []string{}
		for _, d := range a.AlternateDomains {
			if d == domain {
				continue
			}
			newList = append(newList, d)
		}

		// now put the current primary domain into the alt list
		newList = append(newList, a.PrimaryDomain)
		a.AlternateDomains = newList
	}
	a.PrimaryDomain = domain
	return a.writeSettingsFile()
}

func (a *Account) AddAlternateDomain(domain string) error {
	if strings.TrimSpace(domain) == "" {
		return nil
	}

	// make sure this is g2g
	if a.AlternateDomains == nil {
		a.AlternateDomains = []string{}
	}

	if slices.Contains(a.AlternateDomains, domain) {
		return nil
	}

	a.AlternateDomains = append(a.AlternateDomains, domain)
	return a.writeSettingsFile()
}

func (a *Account) RemoveAlternateDomain(domain string) error {
	newList := []string{}
	for _, d := range a.AlternateDomains {
		if d == domain {
			continue
		}
		newList = append(newList, d)
	}
	a.AlternateDomains = newList
	return a.writeSettingsFile()
}

func (a *Account) AddDatabase(host, user, pass, name string) error {
	fmt.Println("host: "+host, "user: "+user, "pass: "+pass, "name: "+name)
	// verify required params are present
	if strings.TrimSpace(host) == "" {
		return errx.Wrapf("invalid param provided", fmt.Errorf("db-host is missing"), ErrorInvalidParam)
	}
	if strings.TrimSpace(user) == "" {
		return errx.Wrapf("invalid param provided", fmt.Errorf("db-user is missing"), ErrorInvalidParam)
	}
	if strings.TrimSpace(pass) == "" {
		return errx.Wrapf("invalid param provided", fmt.Errorf("db-password is missing"), ErrorInvalidParam)
	}
	if strings.TrimSpace(name) == "" {
		return errx.Wrapf("invalid param provided", fmt.Errorf("db-name is missing"), ErrorInvalidParam)
	}

	db := Database{
		Name:     name,
		Password: pass,
		Host:     host,
		Username: user,
	}

	if a.Databases == nil {
		a.Databases = []Database{}
	}
	a.Databases = append(a.Databases, db)

	return a.writeSettingsFile()
}

func (a *Account) SetUptimeURL(url string) error {
	a.UptimeURL = strings.TrimSpace(url)
	return a.writeSettingsFile()
}
