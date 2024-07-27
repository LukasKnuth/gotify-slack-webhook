package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type ButtonElement struct {
	Text    *TextObject
	Url     string
	Confirm bool // TODO add a Emoji/Note?
}

func (be *ButtonElement) Parse(json *gjson.Result) (Skip, error) {
	if url := json.Get("url"); !url.Exists() {
		// Buttons without URLs are allowed, but we don't support them.
		return true, nil
	} else {
		be.Url = url.String()
		// The actual confirmation logic is irrelevant, just that the sender whats the action confirmed first
		be.Confirm = json.Get("confirm").Exists()
		text := json.Get("text")
		text_object, skip, err := ToTextObject(&text)
		if !skip && err == nil {
			be.Text = text_object
			return false, nil
		} else {
			return skip, err
		}
	}
}

func (be *ButtonElement) Render(out *gotify.MarkdownWriter) error {
	if be.Text == nil {
		return nil
	}
	if be.Confirm {
		// TODO we might need to do more if it's a really destructive thing. Maybe render the warning message from the block?
		if err := out.WriteMarkdown("⚠️"); err != nil {
			return err
		}
	}
	if err := out.WriteMarkdown("["); err != nil {
		return err
	}
	if err := be.Text.Render(out); err != nil {
		return err
	}
	if err := out.WriteMarkdownF("](%s)", be.Url); err != nil {
		return err
	}
	if be.Confirm {
		if err := out.WriteMarkdown("⚠️"); err != nil {
			return err
		}
	}
	return nil
}
