package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type ImageElement struct {
	AltText string
	Url     string
}

func (ie *ImageElement) Parse(json *gjson.Result) (Skip, error) {
	if url := json.Get("image_url"); !url.Exists() {
		// Images can use Slack-Hosted images, which we can't support.
		return true, nil
	} else {
		ie.Url = url.String()
		ie.AltText = json.Get("alt_text").String()
		return false, nil
	}
}

func (ie *ImageElement) Render(out *gotify.MarkdownWriter) error {
	if len(ie.Url) > 0 {
		out.WriteMarkdownF("![%s](%s)", ie.AltText, ie.Url)
	}
	return nil
}
