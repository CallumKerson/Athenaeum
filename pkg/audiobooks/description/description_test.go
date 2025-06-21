package description

import (
	"encoding/json"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testRoundTrip tests JSON and TOML marshalling/unmarshalling for a description
func testRoundTrip(t *testing.T, desc Description) {
	t.Helper()

	// JSON round trip
	jsonData, err := json.Marshal(desc)
	require.NoError(t, err)

	var jsonUnmarshalled Description
	err = json.Unmarshal(jsonData, &jsonUnmarshalled)
	require.NoError(t, err)
	assert.Equal(t, desc, jsonUnmarshalled)

	// TOML round trip
	tomlData, err := toml.Marshal(desc)
	require.NoError(t, err)

	var tomlUnmarshalled Description
	err = toml.Unmarshal(tomlData, &tomlUnmarshalled)
	require.NoError(t, err)
	assert.Equal(t, desc, tomlUnmarshalled)
}

func TestDescription_JSONMarshaling(t *testing.T) {
	desc := Description{
		Text:   "A thrilling adventure",
		Format: Markdown,
	}

	// Marshal to JSON
	data, err := json.Marshal(desc)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Description
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Should be equal
	assert.Equal(t, desc, unmarshaled)
}

func TestDescription_TOMLMarshaling(t *testing.T) {
	desc := Description{
		Text:   "A compelling story with complex characters",
		Format: HTML,
	}

	// Marshal to TOML
	data, err := toml.Marshal(desc)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Description
	err = toml.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Should be equal
	assert.Equal(t, desc, unmarshaled)
}

func TestDescription_EmptyValues(t *testing.T) {
	desc := Description{
		Text:   "",
		Format: Plain,
	}

	// JSON round trip
	jsonData, err := json.Marshal(desc)
	require.NoError(t, err)

	var jsonUnmarshaled Description
	err = json.Unmarshal(jsonData, &jsonUnmarshaled)
	require.NoError(t, err)
	assert.Equal(t, desc, jsonUnmarshaled)

	// TOML round trip
	tomlData, err := toml.Marshal(desc)
	require.NoError(t, err)

	var tomlUnmarshaled Description
	err = toml.Unmarshal(tomlData, &tomlUnmarshaled)
	require.NoError(t, err)
	assert.Equal(t, desc, tomlUnmarshaled)
}

func TestDescription_ZeroValue(t *testing.T) {
	var desc Description

	// Zero value should have empty text and Undefined format (default)
	assert.Equal(t, "", desc.Text)
	assert.Equal(t, Undefined, desc.Format)

	// JSON round trip with zero value
	jsonData, err := json.Marshal(desc)
	require.NoError(t, err)

	var jsonUnmarshaled Description
	err = json.Unmarshal(jsonData, &jsonUnmarshaled)
	require.NoError(t, err)
	assert.Equal(t, desc, jsonUnmarshaled)
}

func TestDescription_LongSummary(t *testing.T) {
	longText := `This is a very long description that contains multiple sentences.
It spans multiple lines and includes various formatting elements.
The story follows a protagonist through their journey of discovery
and self-realisation in a fantastical world filled with magic and wonder.
There are complex character relationships, intricate plot twists,
and beautiful descriptive passages that bring the world to life.`

	desc := Description{
		Text:   longText,
		Format: Markdown,
	}

	testRoundTrip(t, desc)
}

func TestDescription_SpecialCharacters(t *testing.T) {
	desc := Description{
		Text:   `Special characters: "quotes", 'apostrophes', & ampersands, <tags>, [brackets], {braces}, and unicode: 🎧📚`,
		Format: HTML,
	}

	testRoundTrip(t, desc)
}
