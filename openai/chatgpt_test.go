package openai

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestClient_Ask(t *testing.T) {
	key := os.Getenv("test_api_key")
	NewClient(key, true, "")
	answer, err := ChatGPT.Ask("今天是周几？")
	assert.NoError(t, err)
	assert.NotEmpty(t, answer)
	t.Log(answer)
}
