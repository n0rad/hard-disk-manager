package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerFilter_Match(t *testing.T) {
	assert.True(t,  HandlerFilter{}.Match(HandlerFilter{}))
	assert.False(t,  HandlerFilter{Type: "disk"}.Match(HandlerFilter{}))
	assert.True(t,  HandlerFilter{Type: "disk"}.Match(HandlerFilter{Type:"disk"}))
	assert.False(t,  HandlerFilter{Type: "disk"}.Match(HandlerFilter{Type:"yopla"}))
	assert.False(t,  HandlerFilter{Type: "disk", FSType: "crypto"}.Match(HandlerFilter{Type:"disk"}))
	assert.False(t,  HandlerFilter{Type: "disk", FSType: "crypto"}.Match(HandlerFilter{Type:"disk", FSType: "crypto2"}))
	assert.True(t,  HandlerFilter{Type: "disk", FSType: "crypto"}.Match(HandlerFilter{Type:"disk", FSType: "crypto"}))
}
