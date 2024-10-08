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
	textType := json.Get("type").String()
	if textType == "plain_text" || textType == "mrkdwn" {
		to.Type = textType
		to.Text = json.Get("text").String()
		return false, nil
	} else {
		return true, nil
	}
}

func (to *TextObject) Render(out *gotify.MarkdownWriter) error {
	if len(to.Text) == 0 {
		return nil
	}
	// Just prints it's content and thats it.
	// Containers are responsible for ending markdown blocks (with newlines)
	switch to.Type {
	case "mrkdwn":
		return out.WriteMarkdown(to.Text)
	default:
		return out.WritePlainText(to.Text)
	}
}
