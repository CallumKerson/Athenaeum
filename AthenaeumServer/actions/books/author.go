package books

import (
	"strings"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func HasAuthor(authors []model.Person, authorQuery string) bool {
	for _, author := range authors {
		if author.String() == authorQuery {
			return true
		}
	}
	return false
}

func ParseAuthors(authorsString string) []model.Person {
	var authors []model.Person
	for _, v := range parseAuthorStrings(authorsString) {
		authors = append(authors, parseAuthor(v))
	}
	return authors
}

func parseAuthor(authorString string) model.Person {
	if strings.HasSuffix(authorString, "Le Guin") {
		return model.Person{GivenNames: strings.TrimSpace(strings.TrimSuffix(authorString, "Le Guin")), FamilyName: "Le Guin"}
	} else {
		sl := strings.Split(authorString, " ")
		lastName := sl[len(sl)-1]
		firstNames := strings.TrimSpace(strings.TrimSuffix(authorString, lastName))
		return model.Person{GivenNames: firstNames, FamilyName: lastName}
	}
}

func parseAuthorStrings(authorsString string) []string {
	noAmpersands := strings.ReplaceAll(authorsString, "&", ",")
	noSemicolons := strings.ReplaceAll(noAmpersands, ";", ",")
	commaSeparated := strings.ReplaceAll(noSemicolons, " and ", ", ")
	authors := strings.Split(commaSeparated, ",")
	for i := range authors {
		authors[i] = strings.TrimSpace(authors[i])
	}
	return authors
}
