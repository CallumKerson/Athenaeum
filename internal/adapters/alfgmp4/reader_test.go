package alfgmp4_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CallumKerson/Athenaeum/internal/adapters/alfgmp4"
	"github.com/CallumKerson/Athenaeum/internal/testing/dataloader"
)

func TestRead(t *testing.T) {
	testReader := alfgmp4.NewMetadataReader()
	metadata, err := testReader.Read(dataloader.GetRootTestdata(t, "Ursula K Le Guin", "Earthsea", "1 A Wizard of Earthsea", "A Wizard of Earthsea.m4b"))

	assert.NoError(t, err)
	assert.Equal(t, 4671000064*time.Nanosecond, metadata.Duration)
}
