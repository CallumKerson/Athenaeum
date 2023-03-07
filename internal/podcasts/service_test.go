package podcasts

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

type testAudiobookClient struct{}

func (c *testAudiobookClient) GetAllAudiobooks(ctx context.Context) ([]audiobooks.Audiobook, error) {
	aWizardOfEarthseaReleaseDate, err := audiobooks.NewReleaseDate("1968-11-01")
	if err != nil {
		return nil, err
	}
	return []audiobooks.Audiobook{{
		Title:       "A Wizard of Earthsea",
		Authors:     []string{"Ursula K. Le Guin"},
		Narrators:   []string{"Kobna Holdbrook-Smith"},
		Path:        "/a-wizard-of-earthsea.m4b",
		FileSize:    200061260,
		MIMEType:    "audio/mp4a-latm",
		ReleaseDate: aWizardOfEarthseaReleaseDate,
		Genres:      []audiobooks.Genre{audiobooks.Childrens, audiobooks.Fantasy},
		Description: &audiobooks.Description{Text: "<p>Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.</p><p>Hungry for power and knowledge, Sparrowhawk tampered with long-held secrets and loosed a terrible shadow upon the world. This is the tale of his testing, how he mastered the mighty words of power, tamed an ancient dragon, and crossed death's threshold to restore the balance.</p>", Format: audiobooks.HTML},
	}}, nil
}

func TestGetFeed(t *testing.T) {
	// given
	expectedTestFeed, err := os.ReadFile(filepath.Join("testdata", "feed1.rss"))
	assert.NoError(t, err)

	testOpts := &FeedOpts{
		Title:       "Audiobooks",
		Description: "Like movies for your mind!",
		Explicit:    true,
		Language:    "EN",
		Author:      "A Person",
		Email:       "person@domain.test",
		Copyright:   "None",
		Link:        "http://www.example-podcast.com/audiobooks",
	}

	svc := NewService("http://www.example-podcast.com/audiobooks/media", testOpts, &testAudiobookClient{}, tlogger.NewTLogger(t))

	feed, err := svc.GetFeed(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, strings.Trim(string(expectedTestFeed), "\n"), feed)
}
