package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type HeaderBlock struct {
	// It's a TextObject but only `type=plain_text` is allowed.
	PlainText string
}

func (hb *HeaderBlock) Parse(block *gjson.Result) (Skip, error) {
	if text_type := block.Get("text.type").String(); text_type == "plain_text" {
		hb.PlainText = block.Get("text.text").String()
		return false, nil
	} else {
		return true, nil
	}
}

func (hb *HeaderBlock) Render(out *gotify.MarkdownWriter) error {
	if len(hb.PlainText) > 0 {
		return out.WriteMarkdownF("## %s\n\n", hb.PlainText)
	} else {
		return nil
	}
}
