package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type DividerBlock struct{}

func (db *DividerBlock) Parse(_ *gjson.Result) (Skip, error) {
	// Just a divider, nothing to parse
	return false, nil
}

func (db *DividerBlock) Render(out *gotify.MarkdownWriter) error {
	return out.WriteMarkdownF("\n\n---\n\n")
}
