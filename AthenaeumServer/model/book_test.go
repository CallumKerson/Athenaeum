package model_test

import (
	"testing"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func TestBookDescriptionNoSeries(t *testing.T) {
	//given
	expectedDescription := "Book Title by Book Author"

	b := model.Book{
		Id:     "001",
		Title:  "Book Title",
		Author: []model.Person{{"Book", "Author"}},
	}

	//when
	actualDescription := b.Description()

	//then
	if actualDescription != expectedDescription {
		t.Errorf("Generated wrong description\nExpected: %s\nWas: %s", expectedDescription, actualDescription)
	}
}

func TestBookDescriptionWithSeries(t *testing.T) {
	//given
	expectedDescription := "An Inbetween Book by Book Author, A Long Running Series Book 2.5"

	b := model.Book{
		Id:     "001",
		Title:  "An Inbetween Book",
		Author: []model.Person{{"Book", "Author"}},
		Series: &model.Series{
			Entry: 2.5,
			Title: "A Long Running Series",
		},
	}

	//when
	actualDescription := b.Description()

	//then
	if actualDescription != expectedDescription {
		t.Errorf("Generated wrong description\nExpected: %s\nWas: %s", expectedDescription, actualDescription)
	}
}

func TestBookFullHTMLDescription(t *testing.T) {
	//given
	expectedDescription := "<h1>An Inbetween Book</h1><h2>By Book Author and Book Co-Author</h2><h4>A Long Running Series Book 2.5</h4>"
	expectedDescription += "<p>This a blurb for the book.</p><p>The blurb is really long, so it is split into multiple paragraphs.</p>"

	b := model.Book{
		Id:      "001",
		Title:   "An Inbetween Book",
		Author:  []model.Person{{"Book", "Author"}, {"Book", "Co-Author"}},
		Summary: "This a blurb for the book.\nThe blurb is really long, so it is split into multiple paragraphs.",
		Series: &model.Series{
			Entry: 2.5,
			Title: "A Long Running Series",
		},
	}

	//when
	actualDescription := b.FullHTMLDescription()

	//then
	if actualDescription != expectedDescription {
		t.Errorf("Generated wrong description\nExpected: %s\nWas: %s", expectedDescription, actualDescription)
	}
}

func TestShouldFormatReleaseDate(t *testing.T) {
	//given
	expected := "2009-10-01"

	b := model.Book{
		ReleaseDateTime: time.Date(2009, 10, 01, 8, 34, 58, 651387237, time.UTC),
	}

	//when
	actual := b.ReleaseDate()

	//then
	if actual != expected {
		t.Errorf("Generated wrong release date\nExpected: %s\nWas: %s", expected, actual)
	}
}

func TestShouldGetAuthorStringForOneAuthor(t *testing.T) {
	//given
	expected := "Book Author"

	b := model.Book{
		Author: []model.Person{{"Book", "Author"}},
	}

	//when
	authorString := b.AuthorString()

	//then
	if expected != authorString {
		t.Errorf("Generated wrong author string\nExpected: %s\nWas: %s", expected, authorString)
	}
}

func TestShouldGetAuthorStringForTwoAuthors(t *testing.T) {
	//given
	expected := "Book Author and Another Writer"

	b := model.Book{
		Author: []model.Person{{"Book", "Author"}, {"Another", "Writer"}},
	}

	//when
	authorString := b.AuthorString()

	//then
	if expected != authorString {
		t.Errorf("Generated wrong author string\nExpected: %s\nWas: %s", expected, authorString)
	}
}

func TestShouldGetAuthorStringForThreeAuthors(t *testing.T) {
	//given
	expected := "Author One, Writer Two and Novelist Three"

	b := model.Book{
		Author: []model.Person{{"Author", "One"}, {"Writer", "Two"}, {"Novelist", "Three"}},
	}

	//when
	authorString := b.AuthorString()

	//then
	if expected != authorString {
		t.Errorf("Generated wrong author string\nExpected: %s\nWas: %s", expected, authorString)
	}
}
