package model

import "time"

type PodcastAuthor struct {
	AuthorName  string
	AuthorEmail string
}

type Category struct {
	MainCategory  string
	SubCategories []string
}

type PodcastOfBooks struct {
	Title           string
	Link            string
	Description     string
	PublicationDate *time.Time
	LastBuildTime   *time.Time
	Author          PodcastAuthor
	ExplicitStatus  bool
	Category        Category
	Items           []PodcastFeedItem
}

type PodcastFeedItem struct {
	GUID              string
	Title             string
	Subtitle          string
	Description       string
	Summary           string
	PublicationDate   time.Time
	DurationInSeconds int64
	Enclosure         Enclosure
}

type Enclosure struct {
	URL    string
	Length int64
}
