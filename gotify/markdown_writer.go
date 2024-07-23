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

func (mdw *MarkdownWriter) NewLine() error {
	_, err := fmt.Fprint(mdw.Writer, "\n")
	return err
}

func (mdw *MarkdownWriter) WriteMarkdown(data string) error {
	_, err := fmt.Fprint(mdw.Writer, data)
	return err
}

func (mdw *MarkdownWriter) WriteMarkdownLn(data string) error {
	_, err := fmt.Fprintln(mdw.Writer, data)
	return err
}

func (mdw *MarkdownWriter) WriteMarkdownF(format string, data ...any) error {
	_, err := fmt.Fprintf(mdw.Writer, format, data...)
	return err
}

func (mdw *MarkdownWriter) WritePlainText(data string) error {
	return mdw.writeEscape(data)
}

func (mdw *MarkdownWriter) WritePlainTextLn(data string) error {
	err := mdw.writeEscape(data)
	if err != nil {
		return err
	} else {
		return mdw.NewLine()
	}
}

func (mdw *MarkdownWriter) WritePlainTextF(format string, data ...any) error {
	formatted := fmt.Sprintf(format, data...)
	return mdw.writeEscape(formatted)
}

func (mdw *MarkdownWriter) writeEscape(text string) error {
	_, err := mdw.escaper.WriteString(mdw.Writer, text)
	return err
}
