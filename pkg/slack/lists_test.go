package slack

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewListView(t *testing.T) {

	view := NewListView()

	jsonView, _ := json.Marshal(view)
	fmt.Print(string(jsonView))
	assert.NotNil(t, view)
	//assert.Equal(t, expectedCommand, command, "they should be equal")

}
