package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type SectionBlock struct {
	Text   string
	Fields []*TextObject
	// TODO Accessory (either image or button - with URL)
}

func (sb *SectionBlock) Parse(block *gjson.Result) (Skip, error) {
	if text := block.Get("text.text"); text.Exists() {
		sb.Text = text.String()
	}
	block.Get("fields.#").ForEach(func(_, value gjson.Result) bool {
		to, skip, err := ToTextObject(&value)
		if err == nil && !skip {
			// TODO we can't really surface the error here. how?
			sb.Fields = append(sb.Fields, to)
		}
		return true
	})
	return false, nil
}

func (sb *SectionBlock) Render(out *gotify.MarkdownWriter) error {
	if len(sb.Text) > 1 {
		err := out.WriteMarkdownLn(sb.Text)
		if err != nil {
			return err
		}
	}
	for _, text_object := range sb.Fields {
		err := out.NewLine()
		if err != nil {
			return err
		}
		err = text_object.Render(out)
		if err != nil {
			return err
		}
	}
	return nil
}
