package legacy

import (
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/tidwall/gjson"
)

type Attachment struct {
	Pretext    string
	AuthorName string
	AuthorLink string
	Title      string
	Text       string
	Footer     string
	Fields     []*AttachmentField
}

func (la *Attachment) hasContent() bool {
	return la.Text != "" || len(la.Fields) > 0
}

func (la *Attachment) Parse(json *gjson.Result) bool {
	if pretext := json.Get("pretext"); pretext.Exists() {
		la.Pretext = pretext.String()
	}
	if author_name := json.Get("author_name"); author_name.Exists() {
		la.AuthorName = author_name.String()
	}
	if author_link := json.Get("author_link"); author_link.Exists() {
		la.AuthorLink = author_link.String()
	}
	if title := json.Get("title"); title.Exists() {
		la.Title = title.String()
	}
	if text := json.Get("text"); text.Exists() {
		la.Text = text.String()
	}
	if footer := json.Get("footer"); footer.Exists() {
		la.Footer = footer.String()
	}
	json.Get("fields").ForEach(func(_, value gjson.Result) bool {
		field := &AttachmentField{}
		skip := field.Parse(&value)
		if !skip {
			la.Fields = append(la.Fields, field)
		}
		return true
	})
	return !la.hasContent()
}

func (la *Attachment) Render(out *gotify.MarkdownWriter) error {
	var err error
	if la.Pretext != "" {
		err := out.WriteMarkdownF("%s\n\n", la.Pretext)
		if err != nil {
			return err
		}
	}
	if la.AuthorName != "" {
		if la.AuthorLink != "" {
			err = out.WriteMarkdownF("[%s](%s)\n\n", la.AuthorName, la.AuthorLink)
		} else {
			err = out.WriteMarkdownF("%s\n\n", la.AuthorName)
		}
		if err != nil {
			return err
		}
	}
	if la.Title != "" {
		err = out.WriteMarkdownF("### %s\n\n", la.Title)
		if err != nil {
			return err
		}
	}
	if la.Text != "" {
		err = out.WriteMarkdownF("%s\n\n", la.Text)
		if err != nil {
			return err
		}
	}
	if len(la.Fields) > 0 {
		for _, field := range la.Fields {
			err = field.Render(out)
			if err != nil {
				return err
			}
		}
		err = out.NewLine()
		if err != nil {
			return err
		}
	}
	if la.Footer != "" {
		err = out.WriteMarkdownF("%s\n\n", la.Footer)
	}
	return err
}
