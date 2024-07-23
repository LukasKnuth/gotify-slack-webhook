package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type SectionBlock struct {
	Text   string
	Fields []string
	// TODO Accessory (either image or button - with URL)
}

func (sb *SectionBlock) Parse(block *gjson.Result) (Skip, error) {
	if text := block.Get("text.text"); text.Exists() {
		sb.Text = text.String()
	}
	block.Get("fields.#.text").ForEach(func(_, value gjson.Result) bool {
		// TODO if type plain_text, we need to escape any Markdown characters
		// TODO escape during parsing OR store markdown/plain and escape during render?
		sb.Fields = append(sb.Fields, value.String())
		return true
	})
	return false, nil
}

func (sb *SectionBlock) Render(out *gotify.MarkdownWriter) error {
	if len(sb.Text) > 1 {
		err := out.WriteMarkdown(sb.Text)
		if err != nil {
			return err
		}
	}
	for _, text := range sb.Fields {
		err := out.WriteMarkdown(text)
		if err != nil {
			return err
		}
	}
	return nil
}
