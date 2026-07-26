package audiobooks

import (
	"errors"
	"fmt"
	"slices"
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
	Erotica
)

var (
	errParsingGenre = errors.New("cannot parse genre")
	// genreName is the display spelling of each genre. It is also the directory
	// name its feed is written under, so changing a value moves a feed.
	genreName = map[Genre]string{
		Literary:   "Literary",
		Mystery:    "Mystery",
		Romance:    "Romance",
		Comedy:     "Comedy",
		Childrens:  "Children's",
		YoungAdult: "Young Adult",
		SciFi:      "Science Fiction",
		Fantasy:    "Fantasy",
		NonFiction: "Non-fiction",
		Biography:  "Biography",
		Historical: "Historical Fiction",
		Thriller:   "Thriller",
		Horror:     "Horror",
		LGBT:       "LGBT+",
		Erotica:    "Erotica",
	}
	// genreValue is every spelling ParseGenre accepts, lowercased. Each one also
	// becomes an alias path for the genre's feed, so entries may be added but
	// never removed: a subscription pointing at one must not start 404ing.
	genreValue = map[string]Genre{
		"literary":           Literary,
		"mystery":            Mystery,
		"romance":            Romance,
		"comedy":             Comedy,
		"children's":         Childrens,
		"children":           Childrens,
		"childrens":          Childrens,
		"young adult":        YoungAdult,
		"youngadult":         YoungAdult,
		"ya":                 YoungAdult,
		"science fiction":    SciFi,
		"sciencefiction":     SciFi,
		"sci-fi":             SciFi,
		"scifi":              SciFi,
		"fantasy":            Fantasy,
		"non-fiction":        NonFiction,
		"nonfiction":         NonFiction,
		"biography":          Biography,
		"historical":         Historical,
		"historical fiction": Historical,
		"historicalfiction":  Historical,
		"thriller":           Thriller,
		"horror":             Horror,
		"lgbt":               LGBT,
		"lgbt+":              LGBT,
		"erotica":            Erotica,
	}
)

// String allows Genre to implement fmt.Stringer.
func (g Genre) String() string {
	return genreName[g]
}

// Aliases returns every string ParseGenre accepts for this genre, sorted so the
// result is stable. Used to emit one feed per spelling a URL might use.
func (g Genre) Aliases() []string {
	aliases := []string{}
	for name, genre := range genreValue {
		if genre == g {
			aliases = append(aliases, name)
		}
	}
	slices.Sort(aliases)
	return aliases
}

// AllGenres returns every defined genre. Map iteration is random, so the result
// is sorted by value — which is declaration order, since the values come from
// iota — to keep the set of genre feeds stable between builds.
func AllGenres() []Genre {
	genres := make([]Genre, 0, len(genreName))
	for genre := range genreName {
		genres = append(genres, genre)
	}
	slices.Sort(genres)
	return genres
}

// Convert a string to a Genre, returns an error if the string is unknown.
func ParseGenre(s string) (Genre, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	genre, ok := genreValue[s]
	if !ok {
		return UndefinedGenre, fmt.Errorf("%w: %q is not a valid genre", errParsingGenre, s)
	}

	return genre, nil
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
