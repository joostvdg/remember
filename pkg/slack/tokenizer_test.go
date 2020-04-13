package slack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenizeSlackCommand_ShouldReturnValidCommandAndSubCommand_2(t *testing.T) {
	text := "view list"
	expectedCommand := "view"
	expectedTokens := []string{"list"}

	command, tokens := tokenizeSlackCommand(text)
	assert.Equal(t, expectedCommand, command, "they should be equal")
	assert.Equal(t, expectedTokens, tokens, "they should be equal")
}

func TestTokenizeSlackCommand_ShouldReturnValidCommandAndSubCommand_1(t *testing.T) {
	text := "view list"
	expectedCommand := "view"
	expectedTokens := []string{"list"}

	command, tokens := tokenizeSlackCommand(text)
	assert.Equal(t, expectedCommand, command, "they should be equal")
	assert.Equal(t, expectedTokens, tokens, "they should be equal")
}
