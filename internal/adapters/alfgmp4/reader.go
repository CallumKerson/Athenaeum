package alfgmp4

import (
	"os"
	"time"

	"github.com/alfg/mp4"

	"github.com/CallumKerson/Athenaeum/pkg/m4b"
)

type M4BMetadataReader struct{}

func NewMetadataReader() *M4BMetadataReader {
	return &M4BMetadataReader{}
}

func (r *M4BMetadataReader) Read(pathToM4BFile string) (*m4b.Metadata, error) {
	metadata := &m4b.Metadata{}
	fInfo, err := os.Stat(pathToM4BFile)
	if err != nil {
		return metadata, err
	}
	file, err := os.Open(pathToM4BFile)
	if err != nil {
		return metadata, err
	}
	defer file.Close()
	info, err := mp4.OpenFromReader(file, fInfo.Size())
	if err != nil {
		return metadata, err
	}

	if info.Moov != nil && info.Moov.Mvhd != nil {
		metadata.Duration = time.Duration(
			(float32(info.Moov.Mvhd.Duration) / float32(info.Moov.Mvhd.Timescale)) * float32(time.Second))
	}
	return metadata, nil
}
