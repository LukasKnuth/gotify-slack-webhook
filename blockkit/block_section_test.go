package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const sectionBlockValidImage = `{
	"type": "section",
	"text": {
		"type": "mrkdwn",
		"text": "Cat quizz"
	},
	"fields": [
		{
			"type": "plain_text",
			"text": "Is it Maggie?"
		},
		{
			"type": "plain_text",
			"text": "Or Jonny?"
		}
	],
	"accessory": {
		"type": "image",
		"image_url": "https://pbs.twimg.com/profile_images/625633822235693056/lNGUneLX_400x400.jpg",
		"alt_text": "cute cat"
	}
}`
const sectionBlockValidBtn = `{
	"type": "section",
	"text": {
		"type": "plain_text",
		"text": "Don't hit the button"
	},
	"fields": [
		{
			"type": "plain_text",
			"text": "Seriously..."
		}
	],
	"accessory": {
		"type": "button",
		"text": {
			"type": "plain_text",
			"text": "Click me ;)"
		},
		"url": "https://explode.everything"
	}
}`
const sectionBlockInvalidFields = `{
	"type": "section",
	"text": {
		"type": "plain_text",
		"text": "This is valid"
	},
	"fields": [
		{
			"type": "invalid",
			"property": "value"
		}
	]
}`
const sectionBlockInvalidAccessory = `{
	"type": "section",
	"text": {
		"type": "mrkdwn",
		"text": "Pick a date"
	},
	"accessory": {
		"type": "datepicker",
		"initial_date": "1990-04-28",
		"placeholder": {
			"type": "plain_text",
			"text": "Select a date",
			"emoji": true
		},
		"action_id": "datepicker-action"
	}
}`

func TestSectionBlockParse(t *testing.T) {
	t.Run("does not fail on empty input", func(t *testing.T) {
		block := &SectionBlock{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips unsupported accessory type", func(t *testing.T) {
		block := &SectionBlock{}
		json := gjson.Parse(sectionBlockInvalidAccessory)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Nil(t, block.Accessory)
		assert.Equal(t, "Pick a date", block.Text.Text)
	})

	t.Run("skips unsupported field type", func(t *testing.T) {
		block := &SectionBlock{}
		json := gjson.Parse(sectionBlockInvalidFields)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Empty(t, block.Fields)
		assert.Equal(t, "This is valid", block.Text.Text)
	})

	t.Run("parses valid input with image", func(t *testing.T) {
		block := &SectionBlock{}
		json := gjson.Parse(sectionBlockValidImage)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "Cat quizz", block.Text.Text)
		assert.Equal(t, "Is it Maggie?", block.Fields[0].Text)
		assert.Equal(t, "Or Jonny?", block.Fields[1].Text)
		assert.Len(t, block.Fields, 2)
		assert.IsType(t, &ImageElement{}, block.Accessory)
	})

	t.Run("parses valid input with button", func(t *testing.T) {
		block := &SectionBlock{}
		json := gjson.Parse(sectionBlockValidBtn)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "Don't hit the button", block.Text.Text)
		assert.Equal(t, "Seriously...", block.Fields[0].Text)
		assert.Len(t, block.Fields, 1)
		assert.IsType(t, &ButtonElement{}, block.Accessory)
	})
}

func TestSectionBlockRender(t *testing.T) {
	t.Run("does not fail with empty struct", func(t *testing.T) {
		block := &SectionBlock{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("renders out valid block", func(t *testing.T) {
		block := &SectionBlock{
			Text: &TextObject{
				Text: "My Text", Type: "mrkdwn",
			},
			Fields: []*TextObject{
				{Text: "left", Type: "plain_text"},
				{Text: "right", Type: "plain_text"},
			},
			Accessory: &ButtonElement{
				Text: &TextObject{Text: "Button", Type: "plain_text"},
				Url:  "http://some.where",
			},
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "My Text\nleft\nright\n[Button](http://some.where)\n", buffer.String())
	})
}
