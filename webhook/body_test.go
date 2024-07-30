package webhook

import (
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/blockkit"
	"github.com/stretchr/testify/assert"
)

const webhookBodyInvalid = `{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Hello, Assistant to the Regional Manager Dwight! *Michael Scott* wants to know where you'd like to take the Paper Company investors to dinner tonight.\n\n *Please select a restaurant:*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "actions",
			"elements": [
				{
					"type": "button",
					"text": {
						"type": "plain_text",
						"text": "Farmhouse",
						"emoji": true
					},
					"value": "click_me_123"
				}
			]
		}
	]
}`
const webhookBodyValid = `{
	"text": "Test notification"
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Hello, Assistant to the Regional Manager Dwight! *Michael Scott* wants to know where you'd like to take the Paper Company investors to dinner tonight.\n\n *Please select a restaurant:*"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "actions",
			"elements": [
				{
					"type": "button",
					"text": {
						"type": "plain_text",
						"text": "Farmhouse",
						"emoji": true
					},
					"value": "click_me_123"
				}
			]
		},
		{
			"type": "image",
			"title": {
				"type": "plain_text",
				"text": "I love tacos",
				"emoji": true
			},
			"image_url": "https://assets3.thrillist.com/v1/image/1682388/size/tl-horizontal_main.jpg",
			"alt_text": "delicious tacos"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Author: K A Applegate",
					"emoji": true
				}
			]
		},
		{
			"type": "input",
			"element": {
				"type": "datepicker",
				"initial_date": "1990-04-28",
				"placeholder": {
					"type": "plain_text",
					"text": "Select a date",
					"emoji": true
				},
				"action_id": "datepicker-action"
			},
			"label": {
				"type": "plain_text",
				"text": "Label",
				"emoji": true
			}
		},
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "This is a header block",
				"emoji": true
			}
		},
		{
			"type": "rich_text",
			"elements": [
				{
					"type": "rich_text_section",
					"elements": [
						{
							"type": "text",
							"text": "Hello there, I am a basic rich text block!"
						}
					]
				}
			]
		},
		{
		  "type": "video",
		  "title": {
		    "type": "plain_text",
		    "text": "Use the Events API to create a dynamic App Home",
		    "emoji": true
		  },
		  "title_url": "https://www.youtube.com/watch?v=8876OZV_Yy0",
		  "description": {
		    "type": "plain_text",
		    "text": "Slack sure is nifty!",
		    "emoji": true
		  },
		  "video_url": "https://www.youtube.com/embed/8876OZV_Yy0?feature=oembed&autoplay=1",
		  "alt_text": "Use the Events API to create a dynamic App Home",
		  "thumbnail_url": "https://i.ytimg.com/vi/8876OZV_Yy0/hqdefault.jpg",
		}
	]
}`

func TestWebhookBodyParse(t *testing.T) {
	t.Run("does not fail for empty input", func(t *testing.T) {
		payload := &WebhookBody{}
		err := payload.Parse([]byte(""))
		assert.Nil(t, err)
		assert.Empty(t, payload.Text)
		assert.Empty(t, payload.Blocks)
	})

	t.Run("skips invalid blocks", func(t *testing.T) {
		payload := &WebhookBody{}
		err := payload.Parse([]byte(webhookBodyInvalid))
		assert.Nil(t, err)
		assert.Empty(t, payload.Text)
		assert.IsType(t, &blockkit.SectionBlock{}, payload.Blocks[0])
		assert.IsType(t, &blockkit.DividerBlock{}, payload.Blocks[1])
		assert.Len(t, payload.Blocks, 2)
	})

	t.Run("parses valid input", func(t *testing.T) {
		payload := &WebhookBody{}
		err := payload.Parse([]byte(webhookBodyValid))
		assert.Nil(t, err)
		assert.Equal(t, "Test notification", payload.Text)
		assert.IsType(t, &blockkit.SectionBlock{}, payload.Blocks[0])
		assert.IsType(t, &blockkit.DividerBlock{}, payload.Blocks[1])
		assert.IsType(t, &blockkit.ImageBlock{}, payload.Blocks[2])
		assert.IsType(t, &blockkit.ContextBlock{}, payload.Blocks[3])
		assert.IsType(t, &blockkit.HeaderBlock{}, payload.Blocks[4])
		assert.IsType(t, &blockkit.VideoBlock{}, payload.Blocks[5])
		assert.Len(t, payload.Blocks, 6)
	})
}

func TestWebhookBodyRender(t *testing.T) {
	t.Run("does not fail for empty struct", func(t *testing.T) {
		payload := &WebhookBody{}
		out, err := payload.Render()
		assert.Nil(t, err)
		assert.Empty(t, out)
	})

	t.Run("renders text and blocks", func(t *testing.T) {
		payload := &WebhookBody{
			Text: "The simple text here",
			Blocks: []blockkit.Block{
				&blockkit.HeaderBlock{PlainText: "Headline"},
				&blockkit.DividerBlock{},
			},
		}
		out, err := payload.Render()
		assert.Nil(t, err)
		assert.Equal(t, "The simple text here\n## Headline\n\n\n\n---\n\n", out)
	})
}
