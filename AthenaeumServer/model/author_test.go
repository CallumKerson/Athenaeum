package model_test

import (
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func TestShouldStringifyAuthor(t *testing.T) {
	//given
	expected := "Ursula K Le Guin"
	a := model.Person{GivenNames: "Ursula K", FamilyName: "Le Guin"}

	//when
	actual := a.String()

	//then
	if actual != expected {
		t.Errorf("Generated wrong author\nExpected: %s\nWas: %s", expected, actual)
	}
}
