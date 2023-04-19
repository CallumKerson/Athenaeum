package service

import "github.com/CallumKerson/loggerrific"

type Option func(s *Service)

func WithPodcastFeedInfo(explicit bool, language, author, email, copyright, imageLink string) Option {
	return func(service *Service) {
		service.feedExplicit = explicit
		service.feedLanguage = language
		service.feedAuthor = author
		service.feedAuthorEmail = email
		service.feedCopyright = copyright
		service.feedImageLink = imageLink
	}
}

func WithHandlePreUnixEpoch(handle bool) Option {
	return func(service *Service) {
		service.handlePreUnixEpochDates = handle
	}
}

func WithLogger(logger loggerrific.Logger) Option {
	return func(s *Service) {
		s.log = logger
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
