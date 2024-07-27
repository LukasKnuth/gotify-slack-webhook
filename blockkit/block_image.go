package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type ImageBlock struct {
	Image *ImageElement
	Title string
}

func (ib *ImageBlock) Parse(json *gjson.Result) (Skip, error) {
	image := &ImageElement{}
	skip, err := image.Parse(json)
	if err != nil || skip {
		return skip, err
	} else {
		ib.Image = image
		ib.Title = json.Get("title.text").String()
		return false, nil
	}
}

func (ib *ImageBlock) Render(out *gotify.MarkdownWriter) error {
	if len(ib.Title) > 0 {
		out.WriteMarkdownF("(Image) %s\n", ib.Title)
	}
	if ib.Image != nil {
		ib.Image.Render(out)
	}
	return nil
}
