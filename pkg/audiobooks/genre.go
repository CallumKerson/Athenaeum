package audiobooks

import (
	"encoding/json"
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
)

var (
	errParsingGenre = errors.New("cannot parse genre")
	genreName       = map[uint8]string{
		1: "Literary",
		2: "Mystery",
		3: "Romance",
		4: "Comedy",
		5: "Children's",
		6: "Young Adult",
		7: "Science Fiction",
		8: "Fantasy",
		9: "Non-fiction",
	}
	genreValue = map[string]uint8{
		"literary":        1,
		"mystery":         2,
		"romance":         3,
		"comedy":          4,
		"children's":      5,
		"children":        5,
		"young adult":     6,
		"ya":              6,
		"science fiction": 7,
		"sci-fi":          7,
		"scifi":           7,
		"fantasy":         8,
		"non-fiction":     9,
		"nonfiction":      9,
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

// MarshalJSON allows compatibility with marshalling JSON.
func (g Genre) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON allows compatibility with unmarshalling JSON.
func (g *Genre) UnmarshalJSON(data []byte) (err error) {
	var genre string

	err = json.Unmarshal(data, &genre)
	if err != nil {
		return err
	}

	if *g, err = ParseGenre(genre); err != nil {
		return err
	}

	return nil
}
