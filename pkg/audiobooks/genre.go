package audiobooks

import (
	"errors"
	"fmt"
	"strings"
)

// Genre - representation of genres of books.
type Genre uint8

const (
	UndefinedGenre Genre = iota
	Literary
	Mystery
	Romance
	Comedy
	Childrens
	YoungAdult
	SciFi
	Fantasy
	NonFiction
	Biography
	Historical
	Thriller
	Horror
	LGBT
)

var (
	errParsingGenre = errors.New("cannot parse genre")
	genreName       = map[uint8]string{
		1:  "Literary",
		2:  "Mystery",
		3:  "Romance",
		4:  "Comedy",
		5:  "Children's",
		6:  "Young Adult",
		7:  "Science Fiction",
		8:  "Fantasy",
		9:  "Non-fiction",
		10: "Biography",
		11: "Historical Fiction",
		12: "Thriller",
		13: "Horror",
		14: "LGBT+",
	}
	genreValue = map[string]uint8{
		"literary":           1,
		"mystery":            2,
		"romance":            3,
		"comedy":             4,
		"children's":         5,
		"children":           5,
		"childrens":          5,
		"young adult":        6,
		"youngadult":         6,
		"ya":                 6,
		"science fiction":    7,
		"sciencefiction":     7,
		"sci-fi":             7,
		"scifi":              7,
		"fantasy":            8,
		"non-fiction":        9,
		"nonfiction":         9,
		"biography":          10,
		"historical":         11,
		"historical fiction": 11,
		"historicalfiction":  11,
		"thriller":           12,
		"horror":             13,
		"lgbt":               14,
		"lgbt+":              14,
	}
)

// String allows Genre to implement fmt.Stringer.
func (g Genre) String() string {
	return genreName[uint8(g)]
}

// Convert a string to a Genre, returns an error if the string is unknown.
func ParseGenre(s string) (Genre, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	value, ok := genreValue[s]
	if !ok {
		return Genre(0), fmt.Errorf("%w: %q is not a valid genre", errParsingGenre, s)
	}

	return Genre(value), nil
}

func (g Genre) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

func (g *Genre) UnmarshalText(data []byte) (err error) {
	if *g, err = ParseGenre(string(data)); err != nil {
		return err
	}
	return nil
}
