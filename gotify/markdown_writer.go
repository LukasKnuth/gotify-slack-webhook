package gotify

import (
	"fmt"
	"io"
	"strings"
)

type MarkdownWriter struct {
	Writer  io.Writer
	escaper *strings.Replacer
}

func Wrap(wrap io.Writer) *MarkdownWriter {
	return &MarkdownWriter{
		Writer:  wrap,
		escaper: strings.NewReplacer("*", "\\*", "_", "\\_", "#", "\\#", "-", "\\-", ">", "\\>"),
	}
}

func (mdw *MarkdownWriter) WriteMarkdown(data string) error {
	_, err := fmt.Fprintln(mdw.Writer, data)
	return err
}

func (mdw *MarkdownWriter) WritePlainText(data string) error {
	_, err := mdw.escaper.WriteString(mdw.Writer, data)
	if err != nil {
		return err
	} else {
		_, err = fmt.Fprint(mdw.Writer, "\n")
		return err
	}
}
