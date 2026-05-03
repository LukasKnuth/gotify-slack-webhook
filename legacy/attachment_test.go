package legacy

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestAttachmentParse(t *testing.T) {
	t.Run("skips for empty input", func(t *testing.T) {
		attachment := &Attachment{}
		json := gjson.Parse("")
		skip := attachment.Parse(&json)
		assert.True(t, skip)
	})

	t.Run("skips if no content", func(t *testing.T) {
		attachment := &Attachment{}
		json := gjson.Parse(`{"author_name":"Tester","title":"Test"}`)
		skip := attachment.Parse(&json)
		assert.True(t, skip)
	})

	t.Run("parses valid input", func(t *testing.T) {
		attachment := &Attachment{}
		json := gjson.Parse(`{
			"pretext":"Pretext",
			"author_name":"Author",
			"author_link":"http://test",
			"title":"Title",
			"text":"Text here",
			"footer":"My foot",
			"fields":[{"title":"f_title","value":"f_val","short":true}]
		}`)
		skip := attachment.Parse(&json)
		assert.False(t, skip)
		assert.Equal(t, "Pretext", attachment.Pretext)
		assert.Equal(t, "Author", attachment.AuthorName)
		assert.Equal(t, "http://test", attachment.AuthorLink)
		assert.Equal(t, "Title", attachment.Title)
		assert.Equal(t, "Text here", attachment.Text)
		assert.Equal(t, "My foot", attachment.Footer)
		assert.Equal(t, "f_title", attachment.Fields[0].Title)
		assert.Equal(t, "f_val", attachment.Fields[0].Value)
	})
}

func TestAttachmentRender(t *testing.T) {
	t.Run("renders nothing for empty struct", func(t *testing.T) {
		attachment := &Attachment{}
		buffer := new(bytes.Buffer)
		err := attachment.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("wraps multiline text in blockquote", func(t *testing.T) {
		attachment := &Attachment{Text: "This\nis\na\ntest!"}
		buffer := new(bytes.Buffer)
		err := attachment.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "> This\n> is\n> a\n> test!\n>\n\n", buffer.String())
	})

	t.Run("strips newlines from most data", func(t *testing.T) {
		attachment := &Attachment{
			Pretext:    "Pre\ntext",
			AuthorName: "Lukas\n",
			Title:      "\nTest\ning",
			Text:       "Text\nhere",
			Footer:     "\nBottom\n\n",
		}
		buffer := new(bytes.Buffer)
		err := attachment.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "Pre text\n\n> Lukas\n>\n> ### Test ing\n>\n> Text\n> here\n>\n> Bottom\n\n", buffer.String())
	})

	t.Run("renders full output", func(t *testing.T) {
		attachment := &Attachment{
			Pretext:    "Pretext",
			AuthorName: "Lukas",
			AuthorLink: "http://test",
			Title:      "Testing",
			Text:       "Text here",
			Footer:     "Bottom",
			Fields: []*AttachmentField{
				{Title: "Title", Value: "Value"},
			},
		}
		buffer := new(bytes.Buffer)
		err := attachment.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "Pretext\n\n> [Lukas](http://test)\n>\n> ### Testing\n>\n> Text here\n>\n> - **Title**: Value\n>\n> Bottom\n\n", buffer.String())
	})
}
