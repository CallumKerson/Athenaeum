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

func WithPodcastFeedInfo(explicit bool, language, author, email, copyright string) Option {
	return func(s *Service) {
		s.feedExplicit = explicit
		s.feedLanguage = language
		s.feedAuthor = author
		s.feedAuthorEmail = email
		s.feedCopyright = copyright
	}
}

type FeedOpts struct {
	Title       string
	Description string
	Link        string
	Explicit    bool
	Language    string
	Author      string
	Email       string
	Copyright   string
}
