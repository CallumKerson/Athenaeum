package model

import "fmt"

type Person struct {
	GivenNames string `json:"givenNames"`
	FamilyName string `json:"familyName"`
}

func (a Person) String() string {
	return fmt.Sprintf("%s %s", a.GivenNames, a.FamilyName)
}
