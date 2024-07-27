package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type SectionBlock struct {
	Text      *TextObject
	Fields    []*TextObject
	Accessory Block
}

func (sb *SectionBlock) Parse(block *gjson.Result) (Skip, error) {
	if text := block.Get("text"); text.Exists() {
		sb.Text = &TextObject{}
		sb.Text.Parse(&text)
	}
	block.Get("fields").ForEach(func(_, value gjson.Result) bool {
		to, skip, err := ToTextObject(&value)
		if err == nil && !skip {
			// TODO we can't really surface the error here. how?
			sb.Fields = append(sb.Fields, to)
		}
		return true
	})
	accessory := block.Get("accessory")
	switch accessory.Get("type").String() {
	case "button":
		sectionAccessoryButton(sb, &accessory)
	case "image":
		sectionAccessoryImage(sb, &accessory)
	}
	return false, nil
}

func sectionAccessoryImage(sb *SectionBlock, block *gjson.Result) {
	elem := &ImageElement{}
	skip, err := elem.Parse(block)
	if err == nil && !skip {
		sb.Accessory = elem
	}
}

func sectionAccessoryButton(sb *SectionBlock, block *gjson.Result) {
	elem := &ButtonElement{}
	skip, err := elem.Parse(block)
	if err == nil && !skip {
		sb.Accessory = elem
	}
}

func (sb *SectionBlock) Render(out *gotify.MarkdownWriter) error {
	if sb.Text != nil {
		err := sb.Text.Render(out)
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
	if sb.Accessory != nil {
		err := sb.Accessory.Render(out)
		if err != nil {
			return err
		}
	}
	return nil
}
