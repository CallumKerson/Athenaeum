package podcast_test

import (
	"testing"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/adapters/podcast"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func TestShouldGenerateFeedWhenExplicitSetToTrue(t *testing.T) {

	//given
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>Podcast Title</title>
    <link>https://podcast.test</link>
    <description>A Test Podcast</description>
    <category>Arts</category>
    <generator>go podcast v1.3.1 (github.com/eduncan911/podcast)</generator>
    <language>en-GB</language>
    <lastBuildDate>Tue, 17 Nov 2009 20:34:58 +0000</lastBuildDate>
    <managingEditor>email@podcast.test (Author Name)</managingEditor>
    <pubDate>Thu, 01 Oct 2009 08:34:58 +0000</pubDate>
    <itunes:author>email@podcast.test (Author Name)</itunes:author>
    <itunes:explicit>yes</itunes:explicit>
    <itunes:category text="Arts">
      <itunes:category text="Books"></itunes:category>
    </itunes:category>
  </channel>
</rss>`

	lastBuildTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	lastPubTime := time.Date(2009, 10, 01, 8, 34, 58, 651387237, time.UTC)

	p := model.PodcastOfBooks{
		Title:           "Podcast Title",
		Link:            "https://podcast.test",
		Description:     "A Test Podcast",
		Author:          model.PodcastAuthor{AuthorName: "Author Name", AuthorEmail: "email@podcast.test"},
		ExplicitStatus:  true,
		LastBuildTime:   &lastBuildTime,
		PublicationDate: &lastPubTime,
		Category:        model.Category{MainCategory: "Arts", SubCategories: []string{"Books"}},
	}

	//when
	feed := podcast.NewPodcastGenerator().String(p)

	//then
	if feed != expected {
		t.Errorf("Wrong podcast generation\nExpected:\n%s\nWas:\n%s", expected, feed)
	}

}

func TestShouldGenerateFeedWhenExplicitSetToFalse(t *testing.T) {

	//given
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>Podcast Title</title>
    <link>https://podcast.test</link>
    <description>A Test Podcast</description>
    <category>Arts</category>
    <generator>go podcast v1.3.1 (github.com/eduncan911/podcast)</generator>
    <language>en-GB</language>
    <lastBuildDate>Tue, 17 Nov 2009 20:34:58 +0000</lastBuildDate>
    <managingEditor>email@podcast.test (Author Name)</managingEditor>
    <pubDate>Thu, 01 Oct 2009 08:34:58 +0000</pubDate>
    <itunes:author>email@podcast.test (Author Name)</itunes:author>
    <itunes:explicit>no</itunes:explicit>
    <itunes:category text="Arts">
      <itunes:category text="Books"></itunes:category>
    </itunes:category>
  </channel>
</rss>`

	lastBuildTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	lastPubTime := time.Date(2009, 10, 01, 8, 34, 58, 651387237, time.UTC)

	p := model.PodcastOfBooks{
		Title:           "Podcast Title",
		Link:            "https://podcast.test",
		Description:     "A Test Podcast",
		Author:          model.PodcastAuthor{AuthorName: "Author Name", AuthorEmail: "email@podcast.test"},
		ExplicitStatus:  false,
		LastBuildTime:   &lastBuildTime,
		PublicationDate: &lastPubTime,
		Category:        model.Category{MainCategory: "Arts", SubCategories: []string{"Books"}},
	}

	//when
	feed := podcast.NewPodcastGenerator().String(p)

	//then
	if feed != expected {
		t.Errorf("Wrong podcast generation\nExpected:\n%s\nWas:\n%s", expected, feed)
	}

}

func TestShouldGenerateFeedWhenProvidedItems(t *testing.T) {

	//given
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>Podcast Title</title>
    <link>https://podcast.test</link>
    <description>A Test Podcast</description>
    <category>Arts</category>
    <generator>go podcast v1.3.1 (github.com/eduncan911/podcast)</generator>
    <language>en-GB</language>
    <lastBuildDate>Tue, 17 Nov 2009 20:34:58 +0000</lastBuildDate>
    <managingEditor>email@podcast.test (Author Name)</managingEditor>
    <pubDate>Thu, 01 Oct 2009 08:34:58 +0000</pubDate>
    <itunes:author>email@podcast.test (Author Name)</itunes:author>
    <itunes:explicit>yes</itunes:explicit>
    <itunes:category text="Arts">
      <itunes:category text="Books"></itunes:category>
    </itunes:category>
    <item>
      <guid>unique1</guid>
      <title>Book Title 1</title>
      <link>https://podcast.test/media/1.m4b</link>
      <description>Book Title 1 by Book Author</description>
      <pubDate>Thu, 01 Oct 2009 08:34:58 +0000</pubDate>
      <enclosure url="https://podcast.test/media/1.m4b" length="6000" type="audio/x-m4a"></enclosure>
      <itunes:author>email@podcast.test (Author Name)</itunes:author>
      <itunes:subtitle>Book Author</itunes:subtitle>
      <itunes:summary><![CDATA[This is a book]]></itunes:summary>
    </item>
    <item>
      <guid>unique2</guid>
      <title>Book Title 2</title>
      <link>https://podcast.test/media/2.m4b</link>
      <description>Book Title 2 by Book Author</description>
      <pubDate>Tue, 01 Dec 2009 08:34:58 +0000</pubDate>
      <enclosure url="https://podcast.test/media/2.m4b" length="6600" type="audio/x-m4a"></enclosure>
      <itunes:author>email@podcast.test (Author Name)</itunes:author>
      <itunes:subtitle>Book Author</itunes:subtitle>
      <itunes:summary><![CDATA[This is another book]]></itunes:summary>
    </item>
  </channel>
</rss>`

	lastBuildTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	lastPubTime := time.Date(2009, 10, 01, 8, 34, 58, 651387237, time.UTC)

	p := model.PodcastOfBooks{
		Title:           "Podcast Title",
		Link:            "https://podcast.test",
		Description:     "A Test Podcast",
		Author:          model.PodcastAuthor{AuthorName: "Author Name", AuthorEmail: "email@podcast.test"},
		ExplicitStatus:  true,
		LastBuildTime:   &lastBuildTime,
		PublicationDate: &lastPubTime,
		Category:        model.Category{MainCategory: "Arts", SubCategories: []string{"Books"}},
	}

	i1 := model.PodcastFeedItem{
		GUID:              "unique1",
		Title:             "Book Title 1",
		Subtitle:          "Book Author",
		Description:       "Book Title 1 by Book Author",
		Summary:           "This is a book",
		PublicationDate:   time.Date(2009, 10, 01, 8, 34, 58, 651387237, time.UTC),
		DurationInSeconds: 600,
		Enclosure: model.Enclosure{
			URL:    "https://podcast.test/media/1.m4b",
			Length: 6000,
		},
	}

	i2 := model.PodcastFeedItem{
		GUID:              "unique2",
		Title:             "Book Title 2",
		Subtitle:          "Book Author",
		Description:       "Book Title 2 by Book Author",
		Summary:           "This is another book",
		PublicationDate:   time.Date(2009, 12, 01, 8, 34, 58, 651387237, time.UTC),
		DurationInSeconds: 660,
		Enclosure: model.Enclosure{
			URL:    "https://podcast.test/media/2.m4b",
			Length: 6600,
		},
	}

	p.Items = []model.PodcastFeedItem{i1, i2}

	//when
	feed := podcast.NewPodcastGenerator().String(p)

	//then
	if feed != expected {
		t.Errorf("Wrong podcast generation\nExpected:\n%s\nWas:\n%s", expected, feed)
	}

}
