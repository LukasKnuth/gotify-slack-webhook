package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type Skip bool
type Block interface {
	Parse(block *gjson.Result) (Skip, error)
	Render(out *gotify.MarkdownWriter) error
}
