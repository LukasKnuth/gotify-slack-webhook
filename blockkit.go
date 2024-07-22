package main

import (
	"fmt"
	"io"

	"github.com/gotify/plugin-api"
	"github.com/tidwall/gjson"
)

func ToMessage() (plugin.Message, error) {
	return plugin.Message{
		Title:   "From Slack",
		Message: "Recieved this from Slack...\n\nSome **Markdown** going _down_ right here!",
		Extras:  map[string]interface{}{"client::display": map[string]string{"contentType": "text/markdown"}},
	}, nil
}

type Skip bool
type Block interface {
	Parse(block *gjson.Result) (Skip, error)
	Render(out io.Writer) error
}

type WebhookBody struct {
	Text   string // TODO both in the same document might not be valid...
	Blocks []Block
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

func parseBlock(block *gjson.Result) (Block, Skip, error) {
	switch block.Get("type").String() {
	case "section":
		section := &SectionBlock{}
		skip, err := section.Parse(block)
		return section, skip, err
	default:
		return nil, true, nil
	}
}

// -------- NEW FILE LAter -----------

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
		sb.Fields = append(sb.Fields, value.String())
		return true
	})
	return false, nil
}

func (sb *SectionBlock) Render(out io.Writer) error {
	if len(sb.Text) > 1 {
		_, err := fmt.Fprintln(out, sb.Text)
		if err != nil {
			return err
		}
	}
	for _, text := range sb.Fields {
		_, err := fmt.Fprintln(out, text)
		if err != nil {
			return err
		}
	}
	return nil
}
