package github

import "fmt"

func (p gitHubProvider) CreateRepository(name string) error {
	fmt.Println("i'm github!!!")
	return nil
}
