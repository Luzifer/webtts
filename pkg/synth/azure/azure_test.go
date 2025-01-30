package azure

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSMLGeneration(t *testing.T) {
	r, err := Provider{}.requestSSML("en-US-ChristopherNeural", "en-US", "excited to try text to speech!")
	require.NoError(t, err)

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, string(data), `<speak version="1.0" xml:lang="en-US"><voice name="en-US-ChristopherNeural">excited to try text to speech!</voice></speak>`)
}
