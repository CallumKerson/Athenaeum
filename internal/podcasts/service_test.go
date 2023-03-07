package podcasts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/loggerrific/tlogger"
)

const (
	expectedTestFeed = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<rss xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\" version=\"2.0\">\n  <channel>\n    <title>Audiobooks</title>\n    <link>http://www.example-podcast.com/my-podcast</link>\n    <copyright></copyright>\n    <language>EN</language>\n    <description>Like movies for your mind!</description>\n    <itunes:author>A Person</itunes:author>\n    <itunes:block>yes</itunes:block>\n    <itunes:explicit>yes</itunes:explicit>\n    <itunes:owner>\n      <itunes:name>A Person</itunes:name>\n      <itunes:email>person@domain.test</itunes:email>\n    </itunes:owner>\n    <item>\n      <title>Episode 1</title>\n      <guid>http://www.example-podcast.com/my-podcast/1/episode-one</guid>\n      <pubDate>Thu, 10 Nov 2022 08:00:00 +0000</pubDate>\n      <itunes:duration>3:50</itunes:duration>\n      <enclosure url=\"http://www.example-podcast.com/my-podcast/1/episode.mp3\" length=\"12312\" type=\"MP3\"></enclosure>\n    </item>\n    <item>\n      <title>Episode 2</title>\n      <guid>http://www.example-podcast.com/my-podcast/2/episode-two</guid>\n      <pubDate>Thu, 10 Nov 2022 11:00:00 +0000</pubDate>\n      <itunes:duration>5:20</itunes:duration>\n      <enclosure url=\"http://www.example-podcast.com/my-podcast/2/episode.mp3\" length=\"46732\" type=\"MP3\"></enclosure>\n    </item>\n  </channel>\n</rss>"
)

func TestGetFeed(t *testing.T) {
	testOpts := &FeedOpts{
		Title:       "Audiobooks",
		Description: "Like movies for your mind!",
		Explicit:    true,
		Language:    "EN",
		Author:      "A Person",
		Email:       "person@domain.test",
		Copyright:   "None",
	}

	svc := NewService(testOpts, tlogger.NewTLogger(t))

	feed, err := svc.GetFeed(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, expectedTestFeed, feed)
}
