package dataloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetRootPath(t *testing.T) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)
	for _, err := os.ReadFile(filepath.Join(dir, "go.mod")); err != nil && len(dir) > 1; {
		dir = filepath.Dir(dir)
		_, err = os.ReadFile(filepath.Join(dir, "go.mod"))
	}
	if len(dir) < 2 {
		t.Fail()
	}
	return dir
}

func GetRootTestdata(t *testing.T, subPaths ...string) string {
	baseTestData := []string{GetRootPath(t), "testdata"}
	baseTestData = append(baseTestData, subPaths...)
	return filepath.Join(baseTestData...)
}
