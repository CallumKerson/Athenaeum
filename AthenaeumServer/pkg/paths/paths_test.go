package paths_test

import (
	"path/filepath"
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/paths"
)

const expectedCombination string = "prefix/suffix"

func TestFileCombinationWithTrailingAndLeadingSlashes(t *testing.T) {
	//when
	actual := filepath.Join("prefix/", "/suffix")

	//then
	if expectedCombination != actual {
		t.Errorf("Generated wrong path\nExpected: %s\nWas: %s", expectedCombination, actual)
	}
}

func TestFileCombinationWithOnlyTrailingSlash(t *testing.T) {
	//when
	actual := filepath.Join("prefix/", "suffix")

	//then
	if expectedCombination != actual {
		t.Errorf("Generated wrong path\nExpected: %s\nWas: %s", expectedCombination, actual)
	}
}

func TestFileCombinationWithOnlyLeadingSlash(t *testing.T) {
	//when
	actual := filepath.Join("prefix", "/suffix")

	//then
	if expectedCombination != actual {
		t.Errorf("Generated wrong path\nExpected: %s\nWas: %s", expectedCombination, actual)
	}
}

func TestFileCombinationWithOnlyNoLeadingOrTrailingSlash(t *testing.T) {

	//when
	actual := filepath.Join("prefix", "suffix")

	//then
	if expectedCombination != actual {
		t.Errorf("Generated wrong path\nExpected: %s\nWas: %s", expectedCombination, actual)
	}
}

func TestCleanPathForEmptyPath(t *testing.T) {
	//given
	expected := ""

	//when
	actual := paths.CleanNameOnly("")

	//then
	if expected != actual {
		t.Errorf("Generated wrong clean name\nExpected to be empty\nWas: %s", actual)
	}
}

func TestCleanPathForSelectCharacters(t *testing.T) {
	//given
	expected := "Fire + BrimstoneMarble (Original)"

	//when
	actual := paths.CleanNameOnly("Fire & Brimstone/Marble [Original]")

	//then
	if expected != actual {
		t.Errorf("Generated wrong clean name\nExpected: %s\nWas: %s", expected, actual)
	}
}
