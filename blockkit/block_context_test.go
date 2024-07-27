package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const contextTestValid = `{
	"type": "context",
	"elements": [
		{
			"type": "mrkdwn",
			"text": "*This* is :smile: markdown"
		},
		{
			"type": "image",
			"image_url": "https://pbs.twimg.com/profile_images/625633822235693056/lNGUneLX_400x400.jpg",
			"alt_text": "cute cat"
		},
		{
			"type": "plain_text",
			"text": "Author: K A Applegate",
			"emoji": true
		}
	]
}`
const contextTestInvalid = `{
	"type": "context",
	"elements": [
		{
			"type": "other",
			"special": "We aint in Dallas anymore"
		},
		{
			"type": "plain_text",
			"text": "still here!"
		}
	]
}`

func TestContextBlockParse(t *testing.T) {
	t.Run("does not crash for empty result", func(t *testing.T) {
		block := &ContextBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses a valid example block", func(t *testing.T) {
		block := &ContextBlock{}
		json := gjson.Parse(contextTestValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)

		assert.IsType(t, &TextObject{}, block.Elements[0])
		assert.IsType(t, &ImageElement{}, block.Elements[1])
		assert.IsType(t, &TextObject{}, block.Elements[2])
		assert.Len(t, block.Elements, 3)
	})

	t.Run("skips invalid elements", func(t *testing.T) {
		block := &ContextBlock{}
		json := gjson.Parse(contextTestInvalid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)

		assert.IsType(t, &TextObject{}, block.Elements[0])
		assert.Len(t, block.Elements, 1)
	})
}

func TestContextBlockRender(t *testing.T) {
	t.Run("renders all elements", func(t *testing.T) {
		block := &ContextBlock{
			Elements: []Block{
				&TextObject{Text: "testing", Type: "plain_text"},
				&ImageElement{Url: "http://test.png", AltText: "alting"},
			},
		}

		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)

		assert.Equal(t, "testing\n![alting](http://test.png)\n\n", buffer.String())
	})
}
