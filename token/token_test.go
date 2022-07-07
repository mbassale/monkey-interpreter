package token_test

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	tok := token.NewToken(token.ASSIGN, '=')
	assert.Equal(t, string(tok.Type), token.ASSIGN)
}
