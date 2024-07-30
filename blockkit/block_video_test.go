package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const videoBlockValid = `{
  "type": "video",
  "title": {
    "type": "plain_text",
    "text": "Use the Events API",
    "emoji": true
  },
  "title_url": "https://www.youtube.com/watch?v=8876OZV_Yy0",
  "description": {
    "type": "plain_text",
    "text": "Slack sure is nifty!",
    "emoji": true
  },
  "video_url": "https://youtu.be/8876OZV_Yy0",
  "alt_text": "create a dynamic App Home",
  "thumbnail_url": "https://i.ytimg.com/8876OZV_Yy0.jpg",
}
`

func TestVideoBlockParse(t *testing.T) {
	t.Run("does not fail on empty input", func(t *testing.T) {
		block := &VideoBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses valid input", func(t *testing.T) {
		block := &VideoBlock{}
		json := gjson.Parse(videoBlockValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "Use the Events API", block.Title)
		assert.Equal(t, "https://youtu.be/8876OZV_Yy0", block.VideoUrl)
		assert.Equal(t, "https://i.ytimg.com/8876OZV_Yy0.jpg", block.ThumbUrl)
		assert.Equal(t, "create a dynamic App Home", block.AltText)
	})
}

func TestVideoBlockRender(t *testing.T) {
	t.Run("does not fail for empty struct", func(t *testing.T) {
		block := &VideoBlock{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("renders a video link with preview image", func(t *testing.T) {
		block := &VideoBlock{
			Title:    "My VLog",
			VideoUrl: "https://you.tube/abc213",
			ThumbUrl: "https://i.you.tube/abc213",
			AltText:  "Watch my VLog",
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "[Video: My VLog ![Watch my VLog](https://i.you.tube/abc213)](https://you.tube/abc213)\n", buffer.String())
	})
}
