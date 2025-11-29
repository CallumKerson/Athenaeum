package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
	"github.com/CallumKerson/Athenaeum/internal/testing/testbooks"
	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
	"github.com/CallumKerson/Athenaeum/pkg/m4b"
)

func TestService_GetAudiobooks(t *testing.T) {
	svc := New(&DummyM4BService{}, dataloader.GetRootTestdata(t))

	books, err := svc.GetAllAudiobooks(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

func TestService_UpdateAudiobooks(t *testing.T) {
	svc := New(&DummyM4BService{}, dataloader.GetRootTestdata(t))

	books, changed, err := svc.ScanForNewAndUpdatedAudiobooks(
		context.TODO(),
		[]audiobooks.Audiobook{testbooks.Audiobooks[0]},
	)
	assert.NoError(t, err)
	assert.True(t, changed)
	assert.ElementsMatch(t, testbooks.Audiobooks, books)
}

type DummyM4BService struct{}

func (s *DummyM4BService) Read(pathToM4BFile string) (*m4b.Metadata, error) {
	return &m4b.Metadata{Duration: time.Nanosecond * 4671000064}, nil
}

func TestFindImagePath(t *testing.T) {
	tests := []struct {
		name          string
		imageFile     string
		expectedFound bool
	}{
		{name: "jpg image", imageFile: "test.jpg", expectedFound: true},
		{name: "jpeg image", imageFile: "test.jpeg", expectedFound: true},
		{name: "png image", imageFile: "test.png", expectedFound: true},
		{name: "no image", imageFile: "", expectedFound: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			basePath := filepath.Join(tempDir, "test")

			// Create image file if expected
			if testCase.imageFile != "" {
				imagePath := filepath.Join(tempDir, testCase.imageFile)
				err := os.WriteFile(imagePath, []byte("fake image content"), 0o644)
				require.NoError(t, err)
			}

			// Test findImagePath
			result := findImagePath(basePath)

			if testCase.expectedFound {
				assert.NotEmpty(t, result)
				assert.Contains(t, result, testCase.imageFile)
			} else {
				assert.Empty(t, result)
			}
		})
	}
}

func TestFindImagePath_PrefersJPG(t *testing.T) {
	// Test that jpg is preferred over png when both exist
	tempDir := t.TempDir()
	basePath := filepath.Join(tempDir, "test")

	// Create both jpg and png
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "test.jpg"), []byte("jpg"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "test.png"), []byte("png"), 0o644))

	result := findImagePath(basePath)

	assert.Equal(t, filepath.Join(tempDir, "test.jpg"), result)
}
