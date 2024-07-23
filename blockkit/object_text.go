package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type TextObject struct {
	Text string
	Type string
}

func ToTextObject(json *gjson.Result) (*TextObject, Skip, error) {
	to := &TextObject{}
	if skip, err := to.Parse(json); skip || err != nil {
		return nil, skip, err
	} else {
		return to, false, nil
	}
}

func (to *TextObject) Parse(json *gjson.Result) (Skip, error) {
	to.Type = json.Get("type").String()
	to.Text = json.Get("text").String()
	return false, nil
}

func (to *TextObject) Render(out *gotify.MarkdownWriter) error {
	if len(to.Text) == 0 {
		return nil
	}
	switch to.Type {
	case "mrkdwn":
		return out.WriteMarkdown(to.Text)
	default:
		return out.WritePlainText(to.Text)
	}
}
