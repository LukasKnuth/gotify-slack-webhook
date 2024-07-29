package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const imageElementValid = `{
	"type": "image",
	"image_url": "http://cat.pics/maggie",
	"alt_text": "cute cat"
}`
const imageElementValidNoAlt = `{
	"type": "image",
	"image_url": "http://cat.pics/maggie"
}`
const imageElementInvalid = `{
	"type": "image",
	"slack_file": {
		"url": "https://slack.com/image.png"
	},
	"alt_text": "alt text"
}`

func TestImageElementParse(t *testing.T) {
	t.Run("does not fail for empty input", func(t *testing.T) {
		block := &ImageElement{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips if image isn't hosted publicly", func(t *testing.T) {
		block := &ImageElement{}
		json := gjson.Parse(imageElementInvalid)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("uses empty text if no alt text is specified", func(t *testing.T) {
		block := &ImageElement{}
		json := gjson.Parse(imageElementValidNoAlt)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "", block.AltText)
		assert.Equal(t, "http://cat.pics/maggie", block.Url)
	})

	t.Run("parses valid image", func(t *testing.T) {
		block := &ImageElement{}
		json := gjson.Parse(imageElementValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "cute cat", block.AltText)
		assert.Equal(t, "http://cat.pics/maggie", block.Url)
	})
}

func TestImageElementRender(t *testing.T) {
	t.Run("renders nothing for empty struct", func(t *testing.T) {
		block := &ImageElement{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("renders image markdown tag", func(t *testing.T) {
		block := &ImageElement{
			AltText: "alting", Url: "http://pic.host/image.png",
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "![alting](http://pic.host/image.png)", buffer.String())
	})
}
