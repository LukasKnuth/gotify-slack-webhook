package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const imageBlockValid = `{
	"type": "image",
	"title": {
		"type": "plain_text",
		"text": "I love tacos",
		"emoji": true
	},
	"image_url": "http://url.to/file.png",
	"alt_text": "delicious tacos"
}`
const imageBlockInvalid = `{
	"type": "image",
	"slack_file": {
		"url": "<insert slack file url here>"
	},
	"alt_text": "inspiration"
}`

func TestImageBlockParse(t *testing.T) {
	t.Run("parses an empty input", func(t *testing.T) {
		block := &ImageBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips if the image isn't publicly available", func(t *testing.T) {
		block := &ImageBlock{}
		json := gjson.Parse(imageBlockInvalid)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses a valid input", func(t *testing.T) {
		block := &ImageBlock{}
		json := gjson.Parse(imageBlockValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "I love tacos", block.Title)
		assert.Equal(t, "delicious tacos", block.Image.AltText)
		assert.Equal(t, "http://url.to/file.png", block.Image.Url)
	})
}

func TestImageBlockRender(t *testing.T) {
	t.Run("does not render title if empty", func(t *testing.T) {
		block := &ImageBlock{
			Title: "",
			Image: &ImageElement{
				AltText: "alting",
				Url:     "http://some.where",
			},
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))

		assert.Nil(t, err)
		assert.Equal(t, "![alting](http://some.where)\n", buffer.String())
	})
	// TODO Repeat this for _all_ tests?
	t.Run("does not fail for empty struct", func(t *testing.T) {
		block := &ImageBlock{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))

		assert.Nil(t, err)
		assert.Len(t, buffer.Bytes(), 0)
	})

	t.Run("renders text and image", func(t *testing.T) {
		block := &ImageBlock{
			Title: "cute puppy",
			Image: &ImageElement{
				AltText: "alting",
				Url:     "http://some.where",
			},
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))

		assert.Nil(t, err)
		assert.Equal(t, "(Image) cute puppy\n![alting](http://some.where)\n", buffer.String())
	})
}
