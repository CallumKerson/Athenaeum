package podcasts

import (
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/logging"
)

type Serializer struct {
	generator PodcastGenerator
	logger    logging.Logger
}

type PodcastGenerator interface {
	String(podcast model.PodcastOfBooks) string
}

func NewPodcastSerializer(generator PodcastGenerator, logger logging.Logger) *Serializer {
	return &Serializer{generator, logger}
}

func (serializer *Serializer) String(podcast model.PodcastOfBooks) string {
	return serializer.generator.String(podcast)
}
