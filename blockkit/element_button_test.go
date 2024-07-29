package blockkit

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const buttonElementValid = `{
	"type": "button",
	"text": {
		"type": "plain_text",
		"text": "Click Me",
		"emoji": true
	},
	"value": "click_me_123",
	"url": "https://google.com",
	"action_id": "button-action"
}`
const buttonElementValidConfirm = `{
	"type": "button",
	"text": {
		"type": "plain_text",
		"text": "click me",
		"emoji": true
	},
	"url": "https://google.com",
	"confirm": {
		"title": {
			"type": "plain_text",
			"text": "destructive"
		},
		"text": {
			"type": "plain_text",
			"text": "are you sure?"
		},
		"confirm": {
			"type": "plain_text",
			"text": "yes"
		},
		"deny": {
			"type": "plain_text",
			"text": "no"
		}
	}
}`
const buttonElementInvalidUrl = `{
	"type": "button",
	"text": {
		"type": "plain_text",
		"text": "Click Me",
		"emoji": true
	},
	"value": "click_me_123",
	"action_id": "button-action"
}`
const buttonElementInvalidText = `{
	"type": "button",
	"text": {
		"type": "invalid",
		"text": "Doesn't matter"
	},
	"url": "https://google.com"
}`

func TestButtonElementParse(t *testing.T) {
	t.Run("does not fail on empty input", func(t *testing.T) {
		block := &ButtonElement{}
		json := gjson.Parse("")
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips buttons without URL", func(t *testing.T) {
		block := &ButtonElement{}
		json := gjson.Parse(buttonElementInvalidUrl)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("skips button if text is invalid", func(t *testing.T) {
		block := &ButtonElement{}
		json := gjson.Parse(buttonElementInvalidText)
		skip, err := block.Parse(&json)
		assert.True(t, bool(skip))
		assert.Nil(t, err)
	})

	t.Run("parses valid button", func(t *testing.T) {
		block := &ButtonElement{}
		json := gjson.Parse(buttonElementValid)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "Click Me", block.Text.Text)
		assert.Equal(t, "https://google.com", block.Url)
		assert.False(t, block.Confirm)
	})

	t.Run("parses confirm block as boolean", func(t *testing.T) {
		block := &ButtonElement{}
		json := gjson.Parse(buttonElementValidConfirm)
		skip, err := block.Parse(&json)
		assert.False(t, bool(skip))
		assert.Nil(t, err)
		assert.Equal(t, "click me", block.Text.Text)
		assert.Equal(t, "https://google.com", block.Url)
		assert.True(t, block.Confirm)
	})
}

func TestButtonElementRender(t *testing.T) {
	t.Run("renders nothing on empty struct", func(t *testing.T) {
		block := &ButtonElement{}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("renders a markdown link", func(t *testing.T) {
		block := &ButtonElement{
			Text: &TextObject{
				Text: "Click Me!", Type: "plain_text",
			},
			Url: "https://some.where",
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "[Click Me!](https://some.where)", buffer.String())
	})

	t.Run("renders info if confirmation was required", func(t *testing.T) {
		block := &ButtonElement{
			Text: &TextObject{
				Text: "Danger!", Type: "plain_text",
			},
			Url:     "https://some.where",
			Confirm: true,
		}
		buffer := new(bytes.Buffer)
		err := block.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "⚠️[Danger!](https://some.where)⚠️", buffer.String())
	})
}
