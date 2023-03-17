package service

import "strings"

type Option func(s *Service)

func WithHost(host string) Option {
	return func(s *Service) {
		s.host = strings.TrimSuffix(host, "/")
	}
}

func WithMediaPath(mediaPath string) Option {
	return func(s *Service) {
		s.mediaPath = strings.Trim(mediaPath, "/")
	}
}

func WithPodcastFeedInfo(explicit bool, language, author, email, copyright, imageLink string) Option {
	return func(service *Service) {
		service.feedExplicit = explicit
		service.feedLanguage = language
		service.feedAuthor = author
		service.feedAuthorEmail = email
		service.feedCopyright = copyright
		service.fedImageLink = imageLink
	}
}

type FeedOpts struct {
	Title       string
	Description string
	Link        string
	ImageLink   string
	Explicit    bool
	Language    string
	Author      string
	Email       string
	Copyright   string
}
