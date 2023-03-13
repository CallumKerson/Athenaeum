package service

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks/description"
)

type testAudiobookClient struct{}

func (c *testAudiobookClient) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	return []audiobooks.Audiobook{
		{
			Title:       "A Wizard of Earthsea",
			Authors:     []string{"Ursula K. Le Guin"},
			Narrators:   []string{"Kobna Holdbrook-Smith"},
			Path:        "/a wizard of earthsea.m4b",
			FileSize:    200061260,
			Duration:    (7 * time.Hour) + (5 * time.Second),
			MIMEType:    "audio/mp4a-latm",
			ReleaseDate: &toml.LocalDate{Year: 1968, Month: 11, Day: 1},
			Genres:      []audiobooks.Genre{audiobooks.Childrens, audiobooks.Fantasy},
			Description: &description.Description{Text: "<p>Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.</p><p>Hungry for power and knowledge, Sparrowhawk tampered with long-held secrets and loosed a terrible shadow upon the world. This is the tale of his testing, how he mastered the mighty words of power, tamed an ancient dragon, and crossed death's threshold to restore the balance.</p>", Format: description.HTML},
			Series:      &audiobooks.Series{Sequence: decimal.NewFromInt(1), Title: "Earthsea"},
		},
		{
			Title:       "This Is How You Lose the Time War",
			Authors:     []string{"Amal El-Mohtar", "Max Gladstone"},
			Narrators:   []string{"Cynthia Farrell", "Emily Woo Zeller"},
			Path:        "/this is how you lose the time war.m4b",
			FileSize:    243930066,
			Duration:    (4 * time.Hour) + (16 * time.Minute) + (7 * time.Second),
			MIMEType:    "audio/mp4a-latm",
			ReleaseDate: &toml.LocalDate{Year: 2019, Month: 07, Day: 16},
			Genres:      []audiobooks.Genre{audiobooks.SciFi},
			Description: &description.Description{Text: "Among the ashes of a dying world, an agent of the Commandant finds a letter. It reads: Burn before reading.\nThus begins an unlikely correspondence between two rival agents hellbent on securing the best possible future for their warring factions. Now, what began as a taunt, a battlefield boast, grows into something more. Something epic. Something romantic. Something that could change the past and the future.\nExcept the discovery of their bond would mean death for each of them. There's still a war going on, after all. And someone has to win that war. That's how war works. Right?", Format: description.Plain},
		}}, nil
}

func TestGetFeed(t *testing.T) {
	// given
	expectedTestFeed, err := os.ReadFile(filepath.Join("testdata", "feed1.rss"))
	assert.NoError(t, err)

	svc := New(&testAudiobookClient{},
		tlogger.NewTLogger(t),
		WithHost("http://www.example-podcast.com/audiobooks/"),
		WithMediaPath("/media/"),
		WithPodcastFeedInfo(true, "EN", "A Person", "person@domain.test", "None"))

	var buf bytes.Buffer

	// when
	err = svc.WriteAllAudiobooksFeed(context.TODO(), &buf)

	// then
	assert.NoError(t, err)
	assert.Equal(t, strings.Trim(string(expectedTestFeed), "\n"), buf.String())
}
