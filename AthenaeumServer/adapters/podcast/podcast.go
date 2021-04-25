package podcast

import (
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"

	pod "github.com/eduncan911/podcast"
)

type PodcastGenerator struct {
}

func NewPodcastGenerator() *PodcastGenerator {
	return &PodcastGenerator{}
}

func (generator *PodcastGenerator) String(podcast model.PodcastOfBooks) string {
	p := pod.New(podcast.Title, podcast.Link, podcast.Description, podcast.PublicationDate, podcast.LastBuildTime)
	p.Language = "en-GB"
	p.AddAuthor(podcast.Author.AuthorName, podcast.Author.AuthorEmail)
	p.AddCategory(podcast.Category.MainCategory, podcast.Category.SubCategories)
	if podcast.ExplicitStatus {
		p.IExplicit = "yes"
	} else {
		p.IExplicit = "no"
	}
	for _, book := range podcast.Items {

		_, _ = p.AddItem(getFeedItem(book))
	}
	return p.String()
}

func getFeedItem(item model.PodcastFeedItem) pod.Item {
	it := pod.Item{GUID: item.GUID, Title: item.Title, Description: item.Description}
	it.ISubtitle = item.Subtitle
	it.AddPubDate(&item.PublicationDate)
	it.AddSummary(item.Summary)
	it.AddEnclosure(item.Enclosure.URL, pod.M4A, item.Enclosure.Length)
	return it
}
