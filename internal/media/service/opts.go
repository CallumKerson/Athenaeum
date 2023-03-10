package service

type Option func(s *Service)

func WithPathToMediaRoot(path string) Option {
	return func(s *Service) {
		s.mediaRoot = path
	}
}
