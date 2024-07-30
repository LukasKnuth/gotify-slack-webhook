package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const textObjectValidMd = `{
	"type": "mrkdwn",
	"text": "A markdown text"
}`
const textObjectInvalidType = `{
	"type": "invalid",
	"text": "This is plain"
}`

func TestTextObjectParse(t *testing.T) {
	t.Run("does not fail for empty input", func(t *testing.T) {
		block := &TextObject{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips invalid text types", func(t *testing.T) {
		block := &TextObject{}
		json := gjson.Parse(textObjectInvalidType)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses valid input", func(t *testing.T) {
		block := &TextObject{}
		json := gjson.Parse(textObjectValidMd)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "mrkdwn", block.Type)
		assert.Equal(t, "A markdown text", block.Text)
	})
}

func TestTextObjectNew(t *testing.T) {
	t.Run("does not return an object if not parsed", func(t *testing.T) {
		json := gjson.Parse(textObjectInvalidType)
		block, skip, err := ToTextObject(&json)
		assert.Nil(t, block)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})
}

func TestTextObjectRender(t *testing.T) {
	t.Run("renders nothing for empty struct", func(t *testing.T) {
		block := &TextObject{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("escapes markdown characters in plain_text", func(t *testing.T) {
		block := &TextObject{
			Type: "plain_text", Text: "**Testing**",
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "\\*\\*Testing\\*\\*", buffer.String())
	})

	t.Run("renders markdown as-is", func(t *testing.T) {
		block := &TextObject{
			Type: "mrkdwn", Text: "**Testing**",
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "**Testing**", buffer.String())
	})
}
