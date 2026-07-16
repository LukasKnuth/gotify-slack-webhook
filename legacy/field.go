package legacy

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type AttachmentField struct {
	Title string
	Value string
}

func (af *AttachmentField) Parse(json *gjson.Result) bool {
	if title := json.Get("title"); title.Exists() {
		af.Title = title.String()
	}
	if value := json.Get("value"); value.Exists() {
		af.Value = value.String()
	}
	return af.Value == ""
}

func (af *AttachmentField) Render(out *gotify.MarkdownWriter) error {
	if af.Value != "" {
		if af.Title != "" {
			return out.WriteMarkdownF("- **%s**: %s\n", af.Title, af.Value)
		} else {
			return out.WriteMarkdownF("- %s\n", af.Value)
		}
	} else {
		return nil
	}
}
