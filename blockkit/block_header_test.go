package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const headerTestValid = `{
	"type": "header",
	"text": {
		"type": "plain_text",
		"text": "This is a header block",
		"emoji": true
	}
}`
const headerTestInvalid = `{
	"type": "header",
	"text": {
		"type": "mrkdwn",
		"text": "Not allowed!"
	}
}`

func TestHeaderBlockParse(t *testing.T) {
	t.Run("does not crash for empty input", func(t *testing.T) {
		block := &HeaderBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses valid test input", func(t *testing.T) {
		block := &HeaderBlock{}
		json := gjson.Parse(headerTestValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)

		assert.Equal(t, "This is a header block", block.PlainText)
	})

	t.Run("skips invalid text type", func(t *testing.T) {
		block := &HeaderBlock{}
		json := gjson.Parse(headerTestInvalid)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})
}

func TestHeaderBlockRender(t *testing.T) {
	t.Run("renders nothing if text is empty", func(t *testing.T) {
		block := &HeaderBlock{
			PlainText: "",
		}

		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Len(t, buffer.Bytes(), 0)
	})

	t.Run("renders header", func(t *testing.T) {
		block := &HeaderBlock{
			PlainText: "Test 123",
		}

		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "## Test 123\n\n", buffer.String())
	})
}
