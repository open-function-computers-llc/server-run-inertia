package account

import (
	"os"
	"reflect"
	"testing"
)

func Test_WhenAnAlternateDomainIsSetAsThePrimaryDomainItIsRemovedFromTheAlternateDomainList(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{"www.dummy.com", "dummy.net", "www.dummy.net"},
	}
	a.SetPrimaryDomain("www.dummy.net")

	expectedAlt := []string{"www.dummy.com", "dummy.net", "dummy.com"}
	if !reflect.DeepEqual(a.AlternateDomains, expectedAlt) {
		t.Errorf("Alternate domains doesn't match expected:\ncurrent:  %s\nexpected: %s", a.AlternateDomains, expectedAlt)
	}
}

func Test_IfYouTryAndAddThePrimaryDomainToTheAltListNothingChanges(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{"www.dummy.com"},
	}
	a.AddAlternateDomain("www.dummy.com")

	if len(a.AlternateDomains) != 1 {
		t.Errorf("The alternate domain list should not have changed")
	}
}

func Test_PassingAnEmptyStringToSettingPrimaryDomainChangesNothing(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{"www.dummy.com"},
	}
	a.SetPrimaryDomain("")
	if a.PrimaryDomain != "dummy.com" {
		t.Error("Primary domain was changed")
	}
}

func Test_PassingAnEmptyStringToAddingAdditionalDomainChangesNothing(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{"www.dummy.com"},
	}

	a.AddAlternateDomain("")
	if len(a.AlternateDomains) != 1 {
		t.Errorf("The alternate domain list should not have changed")
	}
}

func Test_SettingTheCurrentPrimaryDomainToTheSameValueChangesNothing(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{},
	}
	a.SetPrimaryDomain("dummy.com")
	if a.PrimaryDomain != "dummy.com" {
		t.Error("Primary domain was changed")
	}

	if len(a.AlternateDomains) != 0 {
		t.Errorf("The alternate domain list should still be empty")
	}
}

func Test_AddingTheSameValueToTheAlternateDomainsListIgnoresTheDuplicate(t *testing.T) {
	a := Account{
		Name:             "dummy",
		PrimaryDomain:    "dummy.com",
		AlternateDomains: []string{"www.dummy.com"},
	}
	a.AddAlternateDomain("www.dummy.com")

	if len(a.AlternateDomains) != 1 {
		t.Errorf("The alternate domain list should still be a single value")
	}
}

func Test_WorkingWithTheDomainsWillInitialiseTheAlternateDomainsSliceIfItIsNull(t *testing.T) {
	a := Account{
		Name:          "dummy",
		PrimaryDomain: "yay.wtf",
	}
	if a.AlternateDomains != nil {
		t.Error("This should be uninitialized at this point")
	}

	a.AddAlternateDomain("yay.com")
	if len(a.AlternateDomains) != 1 {
		t.Error("This should be initialized at this point")
	}
}

func Test_WeCanLockAnAccount(t *testing.T) {
	a := Account{
		Name:     "whatever",
		IsLocked: false,
	}
	a.Lock()

	if !a.IsLocked {
		t.Error("Account should be locked")
	}

	// calling lock on a locked account isn't a problem
	err := a.Lock()
	if err != nil {
		t.Error("Calling lock on a locked account shouldn't return an error" + err.Error())
	}
}

func Test_AnInvalidAccountDoesNotHaveAValidPathToTheSettingsFile(t *testing.T) {
	a := Account{
		Name: "",
	}
	path, err := a.stateDirectory()
	if path != "" || err != nil {
		t.Error("An account with an invalid name shouldn't have a state directory path")
	}
}

func Test_IfAnAccountWithAVaildNameIsMissingTheStateFilePathItWillBeCreatedAutomatically(t *testing.T) {
	// when the ENVVAR "ACCOUNTS_ROOT" is not set, it'll default to /tmp
	// sanity check
	_, err := os.Stat("/tmp/accountName")
	if !os.IsNotExist(err) {
		t.Error("The testing folder /tmp/accountName already exists")
		return
	}

	a := Account{
		Name: "accountName",
	}
	folderPath, _ := a.stateDirectory()

	_, err = os.Stat("/tmp/accountName")
	if err != nil {
		t.Error("The testing folder /tmp/accountName should have been created")
		return
	}
	os.RemoveAll(folderPath)
}

func Test_WeCanRemoveADomainFromTheAlternateDomainsList(t *testing.T) {
	a := Account{
		Name:             "dummy",
		AlternateDomains: []string{"dummy.com", "dummy.net"},
	}

	a.RemoveAlternateDomain("dummy.com")
	expectedAlt := []string{"dummy.net"}
	if !reflect.DeepEqual(a.AlternateDomains, expectedAlt) {
		t.Errorf("Alternate domains doesn't match expected:\ncurrent:  %s\nexpected: %s", a.AlternateDomains, expectedAlt)
	}
}
