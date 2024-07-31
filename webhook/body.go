package webhook

import (
	"bytes"

	"github.com/lukasknuth/gotify-slack-webhook/blockkit"
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type WebhookBody struct {
	Text   string
	Blocks []blockkit.Block
}

func (wb *WebhookBody) Parse(requestBody []byte) error {
	json := gjson.ParseBytes(requestBody)
	if text := json.Get("text"); text.Exists() {
		wb.Text = text.String()
	}
	blocks := json.Get("blocks")
	for _, block := range blocks.Array() {
		parsed, skip, err := parseBlock(&block)
		if err != nil {
			// TODO this solves my exit early issue on blocks, no?
			return err
		} else if parsed == nil || skip {
			continue
		} else {
			wb.Blocks = append(wb.Blocks, parsed)
		}
	}
	return nil
}

func parseBlock(block *gjson.Result) (blockkit.Block, blockkit.Skip, error) {
	switch block.Get("type").String() {
	case "section":
		section := &blockkit.SectionBlock{}
		skip, err := section.Parse(block)
		return section, skip, err
	case "divider":
		divider := &blockkit.DividerBlock{}
		skip, err := divider.Parse(block)
		return divider, skip, err
	case "image":
		image := &blockkit.ImageBlock{}
		skip, err := image.Parse(block)
		return image, skip, err
	case "context":
		context := &blockkit.ContextBlock{}
		skip, err := context.Parse(block)
		return context, skip, err
	case "header":
		header := &blockkit.HeaderBlock{}
		skip, err := header.Parse(block)
		return header, skip, err
	case "video":
		video := &blockkit.VideoBlock{}
		skip, err := video.Parse(block)
		return video, skip, err
	default:
		return nil, true, nil
	}
}

func (wb *WebhookBody) Render() (string, error) {
	buffer := new(bytes.Buffer)
	out := gotify.Wrap(buffer)

	if len(wb.Text) > 0 {
		err := out.WriteMarkdownLn(wb.Text)
		if err != nil {
			return "", err
		}
	}
	for _, block := range wb.Blocks {
		err := block.Render(out)
		if err != nil {
			return "", err
		}
	}

	result := buffer.String()
	return result, nil
}
