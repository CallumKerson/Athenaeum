package audiobooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudiobooks_Persons(t *testing.T) {
	tests := []struct {
		name             string
		expectsNarrators bool
		audiobook        Audiobook
		expectedPersons  string
	}{
		{name: "Single author", expectsNarrators: false, audiobook: Audiobook{Authors: []string{"Ursula K LeGuin"}}, expectedPersons: "Ursula K LeGuin"},
		{name: "Single narrator", expectsNarrators: true, audiobook: Audiobook{Narrators: []string{"Kobna Holdbrook-Smith"}}, expectedPersons: "Kobna Holdbrook-Smith"},
		{name: "Two authors", expectsNarrators: false, audiobook: Audiobook{Authors: []string{"Amal El-Mohtar", "Max Gladstone"}}, expectedPersons: "Amal El-Mohtar & Max Gladstone"},
		{name: "Two narrators", expectsNarrators: true, audiobook: Audiobook{Narrators: []string{"Cynthia Farrell", "Emily Woo Zeller"}}, expectedPersons: "Cynthia Farrell & Emily Woo Zeller"},
		{name: "Multiple narrators", expectsNarrators: true,
			audiobook: Audiobook{Narrators: []string{
				"Jay Snyder", "Brandon Rubin", "Fred Berman", "Lauren Fortgang", "Roger Clark", "Elizabeth Evans", "Tristan Morris"},
			},
			expectedPersons: "Jay Snyder, Brandon Rubin, Fred Berman, Lauren Fortgang, Roger Clark, Elizabeth Evans & Tristan Morris"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var personString string
			if testCase.expectsNarrators {
				personString = testCase.audiobook.GetNarrator()
			} else {
				personString = testCase.audiobook.GetAuthor()
			}
			assert.Equal(t, testCase.expectedPersons, personString)
		})
	}
}
