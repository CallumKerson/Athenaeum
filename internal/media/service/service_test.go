package service

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

var (
	expectedAudiobooks = []audiobooks.Audiobook{
		{
			Title:       "This Is How You Lose the Time War",
			Authors:     []string{"Amal El-Mohtar", "Max Gladstone"},
			Narrators:   []string{"Cynthia Farrell", "Emily Woo Zeller"},
			Path:        "/Amal El-Mohtar and Max Gladstone/This Is How You Lose the Time War/This Is How You Lose the Time War.m4b",
			FileSize:    145789,
			MIMEType:    "audio/mp4a-latm",
			ReleaseDate: &toml.LocalDate{Year: 2019, Month: 07, Day: 16},
			Genres:      []audiobooks.Genre{audiobooks.SciFi},
			Description: &description.Description{Text: "Among the ashes of a dying world, an agent of the Commandant finds a letter. It reads: Burn before reading.\nThus begins an unlikely correspondence between two rival agents hellbent on securing the best possible future for their warring factions. Now, what began as a taunt, a battlefield boast, grows into something more. Something epic. Something romantic. Something that could change the past and the future.\nExcept the discovery of their bond would mean death for each of them. There's still a war going on, after all. And someone has to win that war. That's how war works. Right?", Format: description.Plain},
		},
		{
			Title:       "A Wizard of Earthsea",
			Authors:     []string{"Ursula K. Le Guin"},
			Narrators:   []string{"Kobna Holdbrook-Smith"},
			Path:        "/Ursula K Le Guin/Earthsea/1 A Wizard of Earthsea/A Wizard of Earthsea.m4b",
			FileSize:    145714,
			MIMEType:    "audio/mp4a-latm",
			ReleaseDate: &toml.LocalDate{Year: 1968, Month: 11, Day: 1},
			Genres:      []audiobooks.Genre{audiobooks.Childrens, audiobooks.Fantasy},
			Description: &description.Description{Text: "<p>Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.</p><p>Hungry for power and knowledge, Sparrowhawk tampered with long-held secrets and loosed a terrible shadow upon the world. This is the tale of his testing, how he mastered the mighty words of power, tamed an ancient dragon, and crossed death's threshold to restore the balance.</p>", Format: description.HTML},
			Series:      &audiobooks.Series{Sequence: decimal.NewFromInt(1), Title: "Earthsea"},
		}}
)

func TestService_GetAudiobooks(t *testing.T) {
	svc := New("testdata", tlogger.NewTLogger(t))

	books, err := svc.GetAudiobooks()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedAudiobooks, books)
}
