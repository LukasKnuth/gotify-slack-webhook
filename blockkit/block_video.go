package blockkit

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type VideoBlock struct {
	Title    string
	VideoUrl string
	ThumbUrl string
	AltText  string
}

func (vb *VideoBlock) Parse(block *gjson.Result) (Skip, error) {
	vb.Title = block.Get("title.text").String()
	vb.AltText = block.Get("alt_text").String()
	vb.ThumbUrl = block.Get("thumbnail_url").String()
	vb.VideoUrl = block.Get("video_url").String()
	return len(vb.Title) == 0 && len(vb.ThumbUrl) == 0 && len(vb.VideoUrl) == 0, nil
}

func (vb *VideoBlock) Render(out *gotify.MarkdownWriter) error {
	if len(vb.VideoUrl) > 0 {
		out.WriteMarkdownF("[Video: %s ![%s](%s)](%s)\n", vb.Title, vb.AltText, vb.ThumbUrl, vb.VideoUrl)
	}
	return nil
}
