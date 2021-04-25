package books_test

import (
	"strings"
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/actions/books"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

var expectedAuthorList = []model.Person{
	{GivenNames: "Author", FamilyName: "One"},
	{GivenNames: "Writer", FamilyName: "Two"},
	{GivenNames: "Novelist", FamilyName: "Three"},
}

var expectedAuthorString = personString(expectedAuthorList)

func TestShoulParseAuthorsWithCommasOnlyCorrectly(t *testing.T) {
	// when
	actual := books.ParseAuthors("Author One, Writer Two, Novelist Three")

	//then
	if expectedAuthorString != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expectedAuthorString, personString(actual))
	}
}

func TestShoulParseAuthorsWithSemicolonsOnlyCorrectly(t *testing.T) {
	// when
	actual := books.ParseAuthors("Author One; Writer Two; Novelist Three")

	//then
	if expectedAuthorString != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expectedAuthorString, personString(actual))
	}
}

func TestShoulParseAuthorsWithAmpersandCorrectly(t *testing.T) {
	// when
	actual := books.ParseAuthors("Author One, Writer Two & Novelist Three")

	//then
	if expectedAuthorString != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expectedAuthorString, personString(actual))
	}
}

func TestShoulParseAuthorsWithAndCorrectly(t *testing.T) {
	// when
	actual := books.ParseAuthors("Author One, Writer Two and Novelist Three")

	//then
	if expectedAuthorString != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expectedAuthorString, personString(actual))
	}
}

func TestShoulParseSingleAuthor(t *testing.T) {
	//given
	expected := personString([]model.Person{
		{GivenNames: "Author Number", FamilyName: "One"},
	})

	// when
	actual := books.ParseAuthors("Author Number One")

	//then
	if expected != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expected, personString(actual))
	}
}

func TestShoulParseUrsulaKLeGuinException(t *testing.T) {
	expected := personString([]model.Person{
		{GivenNames: "Ursula K.", FamilyName: "Le Guin"},
		{GivenNames: "Vonda", FamilyName: "McIntyre"},
	})

	// when
	actual := books.ParseAuthors("Ursula K. Le Guin and Vonda McIntyre")

	//then
	if expected != personString(actual) {
		t.Errorf("Failed to correctly parse author\nExpected: %s\nWas: %s", expected, personString(actual))
	}
}

func personString(people []model.Person) string {
	authorStrings := []string{}
	for _, author := range people {
		authorStrings = append(authorStrings, author.String())
	}
	return strings.Join(authorStrings[:], ",")
}
