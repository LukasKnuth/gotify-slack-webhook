package legacy

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestAttachmentFieldParse(t *testing.T) {
	t.Run("skips entry if input empty", func(t *testing.T) {
		field := &AttachmentField{}
		json := gjson.Parse("")
		skip := field.Parse(&json)
		assert.True(t, skip)
	})

	t.Run("skips entry if no content", func(t *testing.T) {
		field := &AttachmentField{}
		json := gjson.Parse(`{"title":"my title"}`)
		skip := field.Parse(&json)
		assert.True(t, skip)
	})

	t.Run("parses valid input", func(t *testing.T) {
		field := &AttachmentField{}
		json := gjson.Parse(`{"title":"hello","value":"world"}`)
		skip := field.Parse(&json)
		assert.False(t, skip)
		assert.Equal(t, "hello", field.Title)
		assert.Equal(t, "world", field.Value)
	})
}

func TestAttachmentFieldRender(t *testing.T) {
	t.Run("renders nothing for empty struct", func(t *testing.T) {
		field := &AttachmentField{}
		buffer := new(bytes.Buffer)
		err := field.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Empty(t, buffer.Bytes())
	})

	t.Run("renders title if present", func(t *testing.T) {
		field := &AttachmentField{
			Title: "Hello",
			Value: "World",
		}
		buffer := new(bytes.Buffer)
		err := field.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "- **Hello**: World\n", buffer.String())
	})

	t.Run("renders value only if no title", func(t *testing.T) {
		field := &AttachmentField{
			Value: "World",
		}
		buffer := new(bytes.Buffer)
		err := field.Render(gotify.Wrap(buffer))
		assert.Nil(t, err)
		assert.Equal(t, "- World\n", buffer.String())
	})
}
