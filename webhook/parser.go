package webhook

import (
	"github.com/lukasknuth/gotify-slack-webhook/blockkit"
	"github.com/tidwall/gjson"
)

type WebhookBody struct {
	Text   string
	Blocks []blockkit.Block
}

func Parse(json string) (*WebhookBody, error) {
	body := &WebhookBody{}
	result := gjson.Get(json, "text")
	if result.Exists() {
		body.Text = result.String()
	}
	result = gjson.Get(json, "blocks")
	for _, block := range result.Array() {
		parsed, skip, err := parseBlock(&block)
		if err != nil {
			return nil, err
		} else if parsed == nil || skip {
			continue
		} else {
			body.Blocks = append(body.Blocks, parsed)
		}
	}
	return body, nil
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
