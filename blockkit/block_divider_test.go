package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const dividerTestValid = `{
	"type": "divider"
}`

func TestDividerBlockParse(t *testing.T) {
	t.Run("does not crash on empty input", func(t *testing.T) {
		block := &DividerBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses valid example input", func(t *testing.T) {
		block := &DividerBlock{}
		json := gjson.Parse(dividerTestValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
	})
}

func TestDividerBlockRender(t *testing.T) {
	t.Run("renders the element", func(t *testing.T) {
		block := &DividerBlock{}

		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "\n\n---\n\n", buffer.String())
	})
}
