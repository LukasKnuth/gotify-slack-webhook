package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type ContextBlock struct {
	// Can be ImageElement or TextObject
	Elements []Block
}

func (cb *ContextBlock) Parse(block *gjson.Result) (Skip, error) {
	block.Get("elements").ForEach(func(_, value gjson.Result) bool {
		switch value.Get("type").String() {
		case "image":
			contextAppendImage(cb, &value)
		case "mrkdwn":
			contextAppendText(cb, &value)
		case "plain_text":
			contextAppendText(cb, &value)
		}
		return true
	})
	return len(cb.Elements) == 0, nil
}

func contextAppendImage(cb *ContextBlock, block *gjson.Result) {
	elem := &ImageElement{}
	skip, err := elem.Parse(block)
	if err == nil && !skip {
		cb.Elements = append(cb.Elements, elem)
	}
}

func contextAppendText(cb *ContextBlock, block *gjson.Result) {
	elem := &TextObject{}
	skip, err := elem.Parse(block)
	if err == nil && !skip {
		cb.Elements = append(cb.Elements, elem)
	}
}

func (cb *ContextBlock) Render(out *gotify.MarkdownWriter) error {
	for _, element := range cb.Elements {
		err := element.Render(out)
		if err != nil {
			return err
		}
		err = out.NewLine()
		if err != nil {
			return err
		}
	}
	return nil
}
